package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

//控制多个chrome tab页
func ManyTabs() {
	// new browser, first tab
	ctx1, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// ensure the first tab is created
	if err := chromedp.Run(ctx1); err != nil {
		panic(err)
	}

	// same browser, second tab
	ctx2, _ := chromedp.NewContext(ctx1)

	// ensure the second tab is created
	if err := chromedp.Run(ctx2); err != nil {
		panic(err)
	}

	c1 := chromedp.FromContext(ctx1)
	c2 := chromedp.FromContext(ctx2)

	fmt.Printf("Same browser: %t\n", c1.Browser == c2.Browser)
	fmt.Printf("Same tab: %t\n", c1.Target == c2.Target)
}

var chrome context.Context = nil
var chrome_cancel context.CancelFunc = nil
var first_tab_cancel context.CancelFunc = nil

func KillCurr(in string) (string, error) {
	first_tab_cancel()
	chrome_cancel()
	chrome = nil
	return "", nil
}

func screenshot(ctx context.Context, p map[string]string, title string) error {
	if p["shot"] != "1" {
		return nil
	}
	if len(p["shotn"]) > 0 {
		title = p["shotn"]
	} else {
		title = title + time.Now().Format("20060102150405") + ".png"
	}
	var buf []byte
	_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
	if err != nil {
		return err
	}

	width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

	// force viewport emulation
	err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
		WithScreenOrientation(&emulation.ScreenOrientation{
			Type:  emulation.OrientationTypePortraitPrimary,
			Angle: 0,
		}).
		Do(ctx)
	if err != nil {
		return err
	}

	// capture screenshot
	buf, err = page.CaptureScreenshot().
		WithQuality(90).
		WithClip(&page.Viewport{
			X:      contentSize.X,
			Y:      contentSize.Y,
			Width:  contentSize.Width,
			Height: contentSize.Height,
			Scale:  1,
		}).Do(ctx)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(title, buf, 0644); err != nil {
		return err
	}
	return nil
}

//go build --buildmode=plugin -o test_v1.0.so plugins/test.go
func RunTask(in string) (string, error) {
	var p map[string]string
	err := json.Unmarshal([]byte(in), &p)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	var title string
	var sleep time.Duration = 0
	if len(p["sleep"]) > 0 {
		i, _ := strconv.Atoi(p["sleep"])
		sleep = time.Duration(i)
	}
	if chrome == nil { //初始创建chrome窗口
		if len(p["path"]) > 1 {
			chromedp.ExecPath("/opt/apps/cn.google.chrome/files/chrome")
		}
		headless := true
		if p["show"] == "1" {
			headless = false
		}
		options := []chromedp.ExecAllocatorOption{
			chromedp.Flag("headless", headless),
			chromedp.Flag("hide-scrollbars", false),
			chromedp.Flag("mute-audio", false),
			chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36`),
		}
		options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
		//创建chrome窗口
		allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
		chrome = allocCtx
		chrome_cancel = cancel

		//keep a tab
		ctx, cancel := chromedp.NewContext(chrome)
		first_tab_cancel = cancel
		var res string

		err = chromedp.Run(ctx,
			chromedp.Navigate(p["url"]),
			chromedp.WaitVisible(p["wait"], chromedp.ByQuery),
			chromedp.Sleep(sleep*time.Second),
			chromedp.Evaluate(p["js"], &res),
			chromedp.Evaluate("document.title", &title),
			chromedp.ActionFunc(func(ctx context.Context) error {
				return screenshot(ctx, p, title)
			}),
		)
		if err != nil {
			return "", err
		}
		return res, nil

	}

	ctx, cancel := chromedp.NewContext(chrome)
	defer cancel()

	var res string
	err = chromedp.Run(ctx,
		chromedp.Navigate(p["url"]),
		chromedp.WaitVisible(p["wait"], chromedp.ByQuery),
		chromedp.Sleep(sleep*time.Second),
		chromedp.Evaluate(p["js"], &res),
		chromedp.Evaluate("document.title", &title),
		chromedp.ActionFunc(func(ctx context.Context) error {
			return screenshot(ctx, p, title)
		}),
	)
	if err != nil {
		return "", err
	}
	return res, nil
}

func main() {
	p := `
		{
			"show":"0",
			"shot":"1",
			"shotn":"sina.png",
			"path":"",
			"url":"http://www.sina.com.cn",
			"wait":"body",
			"sleep":"0",
			"js":"document.location.href"
		}
	`
	r, e := RunTask(p)
	fmt.Println(r)
	fmt.Println(e)

	KillCurr("")
}

func sample() {
	chromedp.ExecPath("/opt/apps/cn.google.chrome/files/chrome")
	//ctx, cancel := chromedp.NewContext(
	//context.Background(),
	//chromedp.WithLogf(log.Printf),
	//)
	//defer cancel()
	//增加选项，允许chrome窗口显示出来
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
		chromedp.Flag("hide-scrollbars", false),
		chromedp.Flag("mute-audio", false),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36`),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	//创建chrome窗口
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var nodes []*cdp.Node
	var executed *runtime.RemoteObject
	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate("http://81.70.218.13/wiki"),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.Evaluate(`var jq = document.createElement('script'); jq.src = "https://cdn.bootcss.com/jquery/1.4.2/jquery.js"; document.getElementsByTagName('head')[0].appendChild(jq);`, &executed),
		chromedp.Nodes(`.//*[@id="test-editormd-view2"]/blockquote[2]/blockquote[1]/p/a[2]`, &nodes),
		chromedp.Evaluate(`$("body").text()`, &res),
	)
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("window object keys: %v", res)
	fmt.Println(res)
	fmt.Println("get nodes:", len(nodes))
	// print titles
	for _, node := range nodes {
		fmt.Println(node.Children[0].NodeValue, ":", node.AttributeValue("href"))
	}
}
