package article_service

import (
	"image"
	"image/draw"
	"image/jpeg"
	"os"

	"go-gin-example/pkg/qrcode"
	"go-gin-example/pkg/file"
)

type ArticlePoster struct {
    PosterName string
    *Article
    Qr *qrcode.QrCode
}

func NewArticlePoster(posterName string, article *Article, qr *qrcode.QrCode) *ArticlePoster {
	return &ArticlePoster{
		PosterName: posterName,
		Article: article,
		Qr: qr,
	}
}

func GetPosterFlag() string {
	return "poster"
}

func (a *ArticlePoster) CheckMergedImage(path string) bool {
    if file.CheckNotExist(path+a.PosterName) == true {
        return false
    }

    return true
}

func (a *ArticlePoster) OpenMergedImage(path string) (*os.File, error) {
    f, err := file.MustOpen(a.PosterName, path)
    if err != nil {
        return nil, err
    }

    return f, nil
}

type Rect struct {
    Name string
    X0   int
    Y0   int
    X1   int
    Y1   int
}

type Pt struct {
    X int
    Y int
}

type ArticlePosterBg struct {
    Name string
    *ArticlePoster
    *Rect
    *Pt
}

func NewArticlePosterBg(name string, ap *ArticlePoster, rect *Rect, pt *Pt) *ArticlePosterBg {
    return &ArticlePosterBg{
        Name:          name,
        ArticlePoster: ap,
        Rect:          rect,
        Pt:            pt,
    }
}

func (a *ArticlePosterBg) Generate() (string, string, error) {
	// 获取二维码存储路径
	fullPath := qrcode.GetQrCodeFullPath()
	// 生成二维码图像
	fileName, path, err := a.Qr.Encode(fullPath)	
	if err != nil {
		return "", "", err
	}
	// 检查合并后图像（指的是存放合并后的海报）是否存在
	if !a.CheckMergedImage(path) {
		// 若不存在，则生成待合并的图像 mergedF
		mergedF, err := a.OpenMergedImage(path)
		if err != nil {
			return "", "", err
		}

		defer mergedF.Close()

		// 打开事先存放的背景图 bgF
		bgF, err := file.MustOpen(a.Name, path)
		if err != nil {
			return "", "", err
		}

		defer bgF.Close()

		// 打开生成的二维码图像 qrF
		qrF, err := file.MustOpen(fileName, path)

		if err != nil {
            return "", "", err
        }
		defer qrF.Close()
		
		// 解码 bgF 和 qrF 返回 image.Image
		bgImage, err := jpeg.Decode(bgF)
		if err != nil {
            return "", "", err
		}
		
		qrImage, err := jpeg.Decode(qrF)
		if err != nil {
            return "", "", err
		}

		// 创建一个新的 RGBA 图像
		jpg := image.NewRGBA(image.Rect(a.Rect.X0, a.Rect.Y0, a.Rect.X1, a.Rect.Y1))
		// 在 RGBA 图像上绘制 背景图（bgF）
		draw.Draw(jpg, jpg.Bounds(), bgImage, bgImage.Bounds().Min, draw.Over)
		// 在已绘制背景图的 RGBA 图像上，在指定 Point 上绘制二维码图像（qrF）
		draw.Draw(jpg, jpg.Bounds(), qrImage, qrImage.Bounds().Min.Sub(image.Pt(a.Pt.X, a.Pt.Y)), draw.Over)

		// 将绘制好的 RGBA 图像以 JPEG 4：2：0 基线格式写入合并后的图像文件（mergedF）
		jpeg.Encode(mergedF, jpg, nil)
	}

	return fileName, path, nil
}