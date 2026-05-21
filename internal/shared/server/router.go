package server

import (
	"fmt"
	"net/http"
	"path"
	"strings"
)

type node struct {
	path       string
	indices    string
	children   []*node
	paramChild *node

	handler        http.Handler
	pattern        string
	paramChildName string

	catchAllChild     *node
	catchAllChildName string
}

type Router struct {
	get, post, put, delete, patch, head, options *node
	custom map[string]*node
	hosts  map[string]*Router
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) getTree(method string) *node {
	switch method {
	case http.MethodGet:
		return r.get
	case http.MethodPost:
		return r.post
	case http.MethodPut:
		return r.put
	case http.MethodDelete:
		return r.delete
	case http.MethodPatch:
		return r.patch
	case http.MethodHead:
		return r.head
	case http.MethodOptions:
		return r.options
	default:
		if r.custom != nil {
			return r.custom[method]
		}
		return nil
	}
}

func (r *Router) setTree(method string, root *node) {
	switch method {
	case http.MethodGet:
		r.get = root
	case http.MethodPost:
		r.post = root
	case http.MethodPut:
		r.put = root
	case http.MethodDelete:
		r.delete = root
	case http.MethodPatch:
		r.patch = root
	case http.MethodHead:
		r.head = root
	case http.MethodOptions:
		r.options = root
	default:
		if r.custom == nil {
			r.custom = make(map[string]*node)
		}
		r.custom[method] = root
	}
}

func (r *Router) Handle(pattern string, handler http.Handler) {
	method, host, routePath := parsePattern(pattern)
	if method == "" || routePath == "" {
		panic(fmt.Sprintf("invalid pattern: %q", pattern))
	}

	target := r
	if host != "" {
		if r.hosts == nil {
			r.hosts = make(map[string]*Router)
		}
		hr, ok := r.hosts[host]
		if !ok {
			hr = NewRouter()
			r.hosts[host] = hr
		}
		target = hr
	}

	root := target.getTree(method)
	if root == nil {
		root = &node{}
		target.setTree(method, root)
	}
	root.addRoute(routePath, pattern, handler)
}

func (r *Router) HandleFunc(pattern string, handler http.HandlerFunc) {
	r.Handle(pattern, handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	p := req.URL.Path

	if pathNeedsClean(p) {
		cp := path.Clean(p)
		if cp != p {
			u := *req.URL
			u.Path = cp
			http.Redirect(w, req, u.String(), http.StatusMovedPermanently)
			return
		}
		p = cp
	}

	target := r
	if r.hosts != nil {
		host := req.Host
		if i := strings.LastIndexByte(host, ':'); i >= 0 {
			host = host[:i]
		}
		if hr, ok := r.hosts[host]; ok {
			target = hr
		}
	}

	if root := target.getTree(req.Method); root != nil {
		if handler, pat := root.search(p, req); handler != nil {
			req.Pattern = pat
			handler.ServeHTTP(w, req)
			return
		}
	}

	target.serveMiss(w, req, p)
}

func (r *Router) serveMiss(w http.ResponseWriter, req *http.Request, p string) {
	if tsPath := trailingSlashAlternate(p); tsPath != "" {
		if root := r.getTree(req.Method); root != nil {
			if h, _ := root.search(tsPath, nil); h != nil {
				http.Redirect(w, req, tsPath, http.StatusMovedPermanently)
				return
			}
		}
	}

	if allow := r.allowedMethods(req.Method, p); allow != "" {
		w.Header().Set("Allow", allow)
		Error(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	Error(w, http.StatusNotFound, "Not Found")
}

func trailingSlashAlternate(p string) string {
	if len(p) <= 1 {
		return ""
	}
	if p[len(p)-1] == '/' {
		return p[:len(p)-1]
	}
	return p + "/"
}

func pathNeedsClean(p string) bool {
	if len(p) == 0 || p[0] != '/' {
		return true
	}
	for i := 0; i < len(p)-1; i++ {
		if p[i] == '/' && (p[i+1] == '/' || p[i+1] == '.') {
			return true
		}
	}
	return false
}

func (r *Router) allowedMethods(skipMethod, p string) string {
	type entry struct {
		name string
		root *node
	}
	methods := [...]entry{
		{http.MethodDelete, r.delete},
		{http.MethodGet, r.get},
		{http.MethodHead, r.head},
		{http.MethodOptions, r.options},
		{http.MethodPatch, r.patch},
		{http.MethodPost, r.post},
		{http.MethodPut, r.put},
	}

	var buf [9]string
	n := 0
	for _, m := range methods {
		if m.root != nil && m.name != skipMethod {
			if h, _ := m.root.search(p, nil); h != nil {
				buf[n] = m.name
				n++
			}
		}
	}
	for method, root := range r.custom {
		if method != skipMethod && n < len(buf) {
			if h, _ := root.search(p, nil); h != nil {
				buf[n] = method
				n++
			}
		}
	}
	if n == 0 {
		return ""
	}
	return strings.Join(buf[:n], ", ")
}

func parsePattern(pattern string) (method, host, routePath string) {
	i := strings.IndexByte(pattern, ' ')
	if i < 0 {
		return "", "", ""
	}
	method = pattern[:i]
	rest := strings.TrimSpace(pattern[i+1:])

	if rest == "" || rest[0] == '/' {
		return method, "", rest
	}
	slash := strings.IndexByte(rest, '/')
	if slash < 0 {
		return "", "", ""
	}
	return method, rest[:slash], rest[slash:]
}

func longestCommonPrefix(a, b string) int {
	max := len(a)
	if len(b) < max {
		max = len(b)
	}
	for i := 0; i < max; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return max
}

func findByte(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

func (n *node) search(searchPath string, req *http.Request) (http.Handler, string) {
	current := n

	for {
		if pLen := len(current.path); pLen > 0 {
			if len(searchPath) < pLen {
				return nil, ""
			}
			if searchPath[:pLen] != current.path {
				return nil, ""
			}
			searchPath = searchPath[pLen:]
		}

		if len(searchPath) == 0 {
			if current.handler != nil {
				return current.handler, current.pattern
			}
			if current.catchAllChild != nil {
				if req != nil {
					req.SetPathValue(current.catchAllChildName, "")
				}
				return current.catchAllChild.handler, current.catchAllChild.pattern
			}
			return nil, ""
		}

		if i := findByte(current.indices, searchPath[0]); i >= 0 && i < len(current.children) {
			current = current.children[i]
			continue
		}

		if current.paramChild != nil {
			end := strings.IndexByte(searchPath, '/')
			if end != 0 {
				var value string
				if end < 0 {
					value = searchPath
					searchPath = ""
				} else {
					_ = searchPath[end]
					value = searchPath[:end]
					searchPath = searchPath[end:]
				}

				if req != nil {
					req.SetPathValue(current.paramChildName, value)
				}

				if len(searchPath) == 0 {
					if current.paramChild.handler != nil {
						return current.paramChild.handler, current.paramChild.pattern
					}
					if current.paramChild.catchAllChild != nil {
						if req != nil {
							req.SetPathValue(current.paramChild.catchAllChildName, "")
						}
						return current.paramChild.catchAllChild.handler, current.paramChild.catchAllChild.pattern
					}
					return nil, ""
				}

				current = current.paramChild
				continue
			}
		}

		if current.catchAllChild != nil {
			if req != nil {
				req.SetPathValue(current.catchAllChildName, searchPath)
			}
			return current.catchAllChild.handler, current.catchAllChild.pattern
		}

		return nil, ""
	}
}

func (n *node) addRoute(routePath, pattern string, handler http.Handler) {
	if n.path == "" && n.children == nil && n.paramChild == nil && n.handler == nil {
		n.insertChild(routePath, pattern, handler)
		return
	}

	n.addRouteInner(routePath, pattern, handler)
}

func (n *node) addRouteInner(routePath, pattern string, handler http.Handler) {
	if routePath == n.path {
		if n.handler != nil {
			panic(fmt.Sprintf("duplicate route for pattern: %q", pattern))
		}
		n.handler = handler
		n.pattern = pattern
		return
	}

	if idx := strings.IndexByte(routePath, '{'); idx >= 0 && (idx < len(n.path) || n.path == routePath[:idx] || strings.HasPrefix(routePath, n.path)) {
		// fall through to normal prefix handling
	}

	lcp := longestCommonPrefix(routePath, n.path)

	if lcp < len(n.path) {
		child := &node{
			path:              n.path[lcp:],
			handler:           n.handler,
			pattern:           n.pattern,
			children:          n.children,
			indices:           n.indices,
			paramChild:        n.paramChild,
			paramChildName:    n.paramChildName,
			catchAllChild:     n.catchAllChild,
			catchAllChildName: n.catchAllChildName,
		}

		n.path = n.path[:lcp]
		n.handler = nil
		n.pattern = ""
		n.children = []*node{child}
		n.indices = string(child.path[0])
		n.paramChild = nil
		n.paramChildName = ""
		n.catchAllChild = nil
		n.catchAllChildName = ""

		remaining := routePath[lcp:]
		if remaining == "" {
			n.handler = handler
			n.pattern = pattern
			return
		}
		n.addChild(remaining, pattern, handler)
		return
	}

	remaining := routePath[lcp:]
	if remaining == "" {
		if n.handler != nil {
			panic(fmt.Sprintf("duplicate route for pattern: %q", pattern))
		}
		n.handler = handler
		n.pattern = pattern
		return
	}

	n.addChild(remaining, pattern, handler)
}

func (n *node) addChild(routePath, pattern string, handler http.Handler) {
	if routePath[0] == '{' {
		end := strings.IndexByte(routePath, '}')
		if end < 0 {
			panic(fmt.Sprintf("missing closing } in pattern: %q", pattern))
		}
		paramName := routePath[1:end]
		remaining := routePath[end+1:]

		if strings.HasSuffix(paramName, "...") {
			paramName = paramName[:len(paramName)-3]
			if remaining != "" {
				panic(fmt.Sprintf("catch-all must be at end of pattern: %q", pattern))
			}
			if n.catchAllChild != nil {
				panic(fmt.Sprintf("duplicate catch-all for pattern: %q", pattern))
			}
			n.catchAllChild = &node{handler: handler, pattern: pattern}
			n.catchAllChildName = paramName
			return
		}

		if n.paramChild == nil {
			n.paramChild = &node{}
			n.paramChildName = paramName
		}

		if remaining == "" {
			if n.paramChild.handler != nil {
				panic(fmt.Sprintf("duplicate route for pattern: %q", pattern))
			}
			n.paramChild.handler = handler
			n.paramChild.pattern = pattern
			return
		}

		n.paramChild.addChild(remaining, pattern, handler)
		return
	}

	for i := 0; i < len(n.indices); i++ {
		if n.indices[i] == routePath[0] {
			n.children[i].addRouteInner(routePath, pattern, handler)
			return
		}
	}

	child := &node{}
	child.insertChild(routePath, pattern, handler)
	n.children = append(n.children, child)
	n.indices += string(routePath[0])
}

func (n *node) insertChild(routePath, pattern string, handler http.Handler) {
	if idx := strings.IndexByte(routePath, '{'); idx >= 0 {
		n.path = routePath[:idx]

		paramPath := routePath[idx:]
		end := strings.IndexByte(paramPath, '}')
		if end < 0 {
			panic(fmt.Sprintf("missing closing } in pattern: %q", pattern))
		}
		paramName := paramPath[1:end]
		remaining := paramPath[end+1:]

		if strings.HasSuffix(paramName, "...") {
			paramName = paramName[:len(paramName)-3]
			if remaining != "" {
				panic(fmt.Sprintf("catch-all must be at end of pattern: %q", pattern))
			}
			n.catchAllChild = &node{handler: handler, pattern: pattern}
			n.catchAllChildName = paramName
			return
		}

		n.paramChild = &node{}
		n.paramChildName = paramName

		if remaining == "" {
			n.paramChild.handler = handler
			n.paramChild.pattern = pattern
		} else {
			n.paramChild.addChild(remaining, pattern, handler)
		}
		return
	}

	n.path = routePath
	n.handler = handler
	n.pattern = pattern
}
