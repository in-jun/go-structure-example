package query

import "testing"

func TestApplyOptions_Defaults(t *testing.T) {
	c := ApplyOptions(nil)
	if c.Page != 1 {
		t.Errorf("expected default Page 1, got %d", c.Page)
	}
	if c.PageSize != 20 {
		t.Errorf("expected default PageSize 20, got %d", c.PageSize)
	}
}

func TestWithPagination(t *testing.T) {
	c := ApplyOptions([]Option{WithPagination(3, 10)})
	if c.Page != 3 {
		t.Errorf("expected Page 3, got %d", c.Page)
	}
	if c.PageSize != 10 {
		t.Errorf("expected PageSize 10, got %d", c.PageSize)
	}
}

func TestConfig_Offset(t *testing.T) {
	tests := []struct {
		page     int
		pageSize int
		want     int
	}{
		{1, 10, 0},
		{2, 10, 10},
		{3, 20, 40},
		{0, 10, 0},
	}
	for _, tt := range tests {
		c := &Config{Page: tt.page, PageSize: tt.pageSize}
		if got := c.Offset(); got != tt.want {
			t.Errorf("Offset(page=%d, pageSize=%d) = %d, want %d", tt.page, tt.pageSize, got, tt.want)
		}
	}
}
