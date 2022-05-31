# SecretImageSharing

2022春 网络信息安全 上机

## Introduction

根据(r,n)门限方案，将一张图片分割成n个图片，取其中任意r张图片可以还原出原图片。

原理基本就是Shamir秘钥分享算法，通过拉格朗日插值法还原多项式系数。

## Quick Start

```bash
Usage of SecretImageSharing.exe:
  -n int
        N (default 5)
  -t int
        T (default 4)
  -tag
        time tag on decrypted image
```

```bash
> SecretImageSharing.exe C:\xxx\xxx.jpg
loaded jpeg image from C:\xxx\xxx.jpg
done

> SecretImageSharing.exe 1 3 4 5
loaded bmp image from image_1.bmp
loaded bmp image from image_3.bmp
loaded bmp image from image_4.bmp
loaded bmp image from image_5.bmp
done
```

原图：
![a.jpeg](https://s2.loli.net/2022/05/31/dfKySBClQGqkLWU.jpg)

黑白化原图：
![a.bmp](https://i.ibb.co/cCDd9wP/image-origin.png)

image_1.bmp：
[![XUHLaS.md.png](https://iili.io/XUHLaS.md.png)](https://freeimage.host/i/XUHLaS)

image_2.bmp：
[![XUHsF2.md.png](https://iili.io/XUHsF2.md.png)](https://freeimage.host/i/XUHsF2)

image_3.bmp：
[![XUHP6l.md.png](https://iili.io/XUHP6l.md.png)](https://freeimage.host/i/XUHP6l)

image_4.bmp：
[![XUH6G4.md.png](https://iili.io/XUH6G4.md.png)](https://freeimage.host/i/XUH6G4)

image_5.bmp：
[![XUH4nf.md.png](https://iili.io/XUH4nf.md.png)](https://freeimage.host/i/XUH4nf)

可以发现低次的图片加密效果并不好，所以可以尽量选取高次的图片。

## References

[https://blog.csdn.net/z784561257/article/details/84424848](https://blog.csdn.net/z784561257/article/details/84424848)
LIN, Pei-Yu; CHAN, Chi-Shiang. Invertible secret image sharing with steganography. Pattern Recognition Letters, 2010, 31.13: 1887-1893.
