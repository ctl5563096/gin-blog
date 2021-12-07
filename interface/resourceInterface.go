package _interface

import (
	"fmt"
	"gin-blog/common"
	"gin-blog/models/blog"
)

type Factory struct {
}

// Product 产品接口
type Product interface {
	GetList(page int,pageSize int) []common.ReturnGetList
}

type ArticleProduct struct {
}

type PhotosProduct struct {
}


func (p1 ArticleProduct) GetList(page int,pageSize int) []common.ReturnGetList {
	var r []common.ReturnGetList
	blog.GetChooseList(common.GetListParams{Page: page ,PageSize: pageSize},&r)
	return r
}

func (p2 PhotosProduct) GetList(page int,pageSize int) []common.ReturnGetList{
	var r []common.ReturnGetList
	fmt.Println(page + pageSize)
	return r
}

// Generate 实现接口
func (f Factory) Generate(resourceType int) Product {
	switch resourceType {
		case 1:
			return ArticleProduct{}
		case 2:
			return PhotosProduct{}
		default:
			return nil
	}
}