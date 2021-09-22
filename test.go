package main

import (
	"fmt"
)

type Info struct {
	sex int
	name string
	age int
	address string
}

type User struct{
	like string
	Info
}

type Admin struct {
	unlike string
	Info
}

type Books struct {
	title string
	author string
	subject string
	book_id int
}

type Person struct {
	Firstname string
	Lastname string
	Age uint8
}

func main()  {
	//user:= User{}
	//user.sex=0
	//user.address="广州市"
	//user.like="游戏"
	//fmt.Println(user)
	//
	//
	//admin:= Admin{Info:Info{sex:1}}//还可以这样声明一些属性值,因为Info是结构体,匿名,所以需要这样声明
	//admin.address="广州市"
	//admin.unlike="游戏"
	//fmt.Println(admin)
	//
	//fmt.Println(admin.Info.sex)

	//var Book1 Books        /* 声明 Book1 为 Books 类型 */
	//var Book2 Books        /* 声明 Book2 为 Books 类型 */
	//
	///* book 1 描述 */
	//Book1.title = "Go 语言"
	//Book1.author = "www.runoob.com"
	//Book1.subject = "Go 语言教程"
	//Book1.book_id = 6495407
	//
	///* book 2 描述 */
	//Book2.title = "Python 教程"
	//Book2.author = "www.runoob.com"
	//Book2.subject = "Python 语言教程"
	//Book2.book_id = 6495700
	//
	///* 打印 Book1 信息 */
	//printBook(&Book1)
	//printBook(&Book1)
	//
	///* 打印 Book2 信息 */
	////printBook(Book2)

	//p := Person{"sun", "xingfang", 30}
	//不一致的情况
	//p.show() // sun 修改前
	//p.setFirstName("tom")   // 值方法
	//p.show2()
	//p.show() // sun, 未变化
	//p.setFirstName2("tom")  // 指针方法
	//p.show() // tom 修改后的tom
}

func printBook( book *Books ) {
	book.book_id = 1232131
	fmt.Printf( "Book title : %s\n", book.title)
	fmt.Printf( "Book author : %s\n", book.author)
	fmt.Printf( "Book subject : %s\n", book.subject)
	fmt.Printf( "Book book_id : %d\n", book.book_id)
}

// 值方法
func (p Person) show() {
	fmt.Println(p.Firstname)
}
// 指针方法
func (p *Person) show2() {
	fmt.Println(p.Firstname)
}

// 值方法
func (p Person) setFirstName(name string) {
	p.Firstname = name
}
// 指针方法
func (p *Person) setFirstName2(name string) {
	p.Firstname = name
}