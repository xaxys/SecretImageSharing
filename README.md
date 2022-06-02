# SecretImageSharing

2022春 网络信息安全 上机

## Introduction

根据(r,n)门限方案，将一张图片分割成n个图片，取其中任意r张图片可以还原出原图片。

原理基本就是Shamir秘钥分享算法，通过拉格朗日插值法还原多项式系数。

- 将图片运用黑白滤镜，变成灰度图片

- ~~使用某些算法，生成一个置换序列，按此序列置换像素点~~
  
  （因为这步置换本身就是一步加密了，在此基础上展现秘密图片分享并不直观，故忽略此步操作，仅仅顺序的置换像素）

- 取原图中连续的r个像素点，作为多项式系数。

  （此处灰度范围为0-255，而其中最大质数为251，多项式应对251取模。这导致251-255的值不能使用，因此提出以下两种方案）

  - 将251-255灰度的像素有损处理为250灰度

  - 使用扩展位，用两位像素点表示一个251-255的灰度值。如：253表示为250+3，读到灰度为250的像素时，再去读下一个像素点，如果为3，则合并为253。

  （此处还有一个问题，解码时计算拉格朗日多项式系数时，高次系数为0时应当被舍弃，所以灰度0也不应当被使用，故在编码时对所有灰度值都加1）

- 运用上述多项式可生成n个点(x,y)，此处x作为图片序号，y作为黑白像素的Y值

- 将上述n个点分别插入n张生成的图片中，最好可以使用伪装图片，把这个像素作为噪声插入图片中，但本程序不做伪装，只是连续的插入像素。

- 解码时，用拉格朗日插值法，还原多项式系数。然后做上面的逆过程即可。

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

数量不足的解码：
[![XUnfgs.md.png](https://iili.io/XUnfgs.md.png)](https://freeimage.host/i/XUnfgs)

数量足够的解码：与原图一致

可以发现低次的图片加密效果并不好，所以可以多生成一些然后尽量选取高次的图片。或将加密的像素作为噪声插入伪装的图片中。

## References

[https://blog.csdn.net/z784561257/article/details/84424848](https://blog.csdn.net/z784561257/article/details/84424848)

LIN, Pei-Yu; CHAN, Chi-Shiang. Invertible secret image sharing with steganography. Pattern Recognition Letters, 2010, 31.13: 1887-1893.
