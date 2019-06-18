package page

import (
	"github.com/solowu/goii/mod/util/gconv"
	"math"
	"net/url"
)

// 分页对象
type Page struct {
	Url         *url.URL // 当前页面的URL对象
	UrlTemplate string   // URL生成规则，内部可使用{.page}变量指定页码
	TotalSize   int64    // 总共数据条数
	TotalPage   int64    // 总页数
	CurrentPage int64    // 当前页码
	QueryName   string   // 分页参数名称(GET参数)
	PageSize    int64    //每页条数
	Offset      int64    //偏移数量
	Form        int64
	To          int64
	_first      *pageItem
	_last       *pageItem
	_prev       *pageItem
	_next       *pageItem
	_list       []*pageItem
}

type pageItem struct {
	Number   int64
	Url      string
	Current  bool
	Disabled bool
}

// 创建一个分页对象，输入参数分别为：
// 总数量、每页数量、当前页码、当前的URL(URI+QUERY)
func New(totalSize, pageSize, currentPage int64, Url string, queryName ...string) *Page {
	u, _ := url.Parse(Url)
	page := &Page{
		Url:         u,
		TotalSize:   totalSize,
		PageSize:    pageSize,
		CurrentPage: currentPage,
		TotalPage:   int64(math.Ceil(float64(totalSize) / float64(pageSize))),
		QueryName:   "page",
	}

	if currentPage <= 1 {
		page.CurrentPage = 1
	}
	if currentPage > page.TotalPage {
		page.CurrentPage = page.TotalPage
	}

	page.Offset = (page.CurrentPage - 1) * pageSize

	if len(queryName) > 0 {
		page.QueryName = queryName[0]
	}

	page.Form = page.Offset + 1
	page.To = page.Offset + page.PageSize
	if page.To > totalSize {
		page.To = totalSize
	}
	return page
}

// 为指定的页面返回地址值
func (page *Page) getUrl(pageNo int64) string {
	// 复制一个URL对象
	u := *page.Url
	values := page.Url.Query()
	if pageNo == 1 {
		values.Del(page.QueryName)
	} else {
		values.Set(page.QueryName, gconv.String(pageNo))
	}
	u.RawQuery = values.Encode()
	return u.String()
}

//获取首页
func (page *Page) First() *pageItem {
	if page._first != nil {
		return page._first
	}
	if page.CurrentPage <= 1 {
		page._first = &pageItem{
			Disabled: true,
		}

	} else {
		page._first = &pageItem{
			Number:  1,
			Url:     page.getUrl(1),
			Current: page.CurrentPage == 1,
		}
	}
	return page._first
}

//尾页
func (page *Page) Last() *pageItem {
	if page._last != nil {
		return page._last
	}
	if page.CurrentPage >= page.TotalPage {
		page._last = &pageItem{
			Disabled: true,
		}
	} else {
		page._last = &pageItem{
			Number: page.TotalPage,
			Url:    page.getUrl(page.TotalPage),
		}
	}
	return page._last
}

//上一页
func (page *Page) Prev() *pageItem {
	if page._prev != nil {
		return page._prev
	}
	if page.CurrentPage <= 1 {
		page._prev = &pageItem{
			Disabled: true,
		}
	} else {
		no := page.CurrentPage - 1
		page._prev = &pageItem{
			Number: no,
			Url:    page.getUrl(no),
		}
	}
	return page._prev
}

//下一页
func (page *Page) Next() *pageItem {
	if page._next != nil {
		return page._next
	}
	if page.CurrentPage >= page.TotalPage {
		page._next = &pageItem{
			Disabled: true,
		}
	} else {
		no := page.CurrentPage + 1
		page._next = &pageItem{
			Number: no,
			Url:    page.getUrl(no),
		}
	}
	return page._next
}

// 获得分页条列表内容
func (page *Page) List(num ...int64) []*pageItem {

	if len(page._list) > 0 {
		return page._list
	}

	var number int64 = 5
	if len(num) > 0 {
		number = num[0]
	}
	begin := int64(math.Ceil(float64(number) / 2))
	/*
		protected function getPageRange()
	    {
	        $currentPage = $this->pagination->getPage();
	        $pageCount = $this->pagination->getPageCount();

	        $beginPage = max(0, $currentPage - (int) ($this->maxButtonCount / 2));
	        if (($endPage = $beginPage + $this->maxButtonCount - 1) >= $pageCount) {
	            $endPage = $pageCount - 1;
	            $beginPage = max(0, $endPage - $this->maxButtonCount + 1);
	        }

	        return [$beginPage, $endPage];
	    }
	*/
	begin = page.CurrentPage - begin
	if begin < 0 {
		begin = 0
	}
	endPage := begin + number - 1
	if endPage >= page.TotalPage {
		endPage = page.TotalPage - 1
		begin = int64(math.Max(0.0, float64(endPage)-float64(number)))
	}
	for begin <= endPage {
		begin++
		item := &pageItem{
			Number: begin,
		}
		if page.CurrentPage == begin {
			item.Current = true
		} else {
			item.Url = page.getUrl(begin)
		}
		page._list = append(page._list, item)

	}

	return page._list
}
