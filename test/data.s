00000000  B80000            mov ax,0x0
00000003  B83412            mov ax,0x1234
00000006  B90100            mov cx,0x1
00000009  BA1000            mov dx,0x10
0000000C  BB0001            mov bx,0x100
0000000F  BC0010            mov sp,0x1000
00000012  BDFF00            mov bp,0xff
00000015  BE00FF            mov si,0xff00
00000018  BFFECA            mov di,0xcafe
0000001B  B000              mov al,0x0
0000001D  B101              mov cl,0x1
0000001F  B210              mov dl,0x10
00000021  B311              mov bl,0x11
00000023  B412              mov ah,0x12
00000025  B5FF              mov ch,0xff
00000027  B6EE              mov dh,0xee
00000029  B7CA              mov bh,0xca
0000002B  8800              mov [bx+si],al
0000002D  8900              mov [bx+si],ax
0000002F  8A00              mov al,[bx+si]
00000031  8B00              mov ax,[bx+si]
00000033  88063412          mov [0x1234],al
00000037  89063412          mov [0x1234],ax
0000003B  8A063412          mov al,[0x1234]
0000003F  8B063412          mov ax,[0x1234]
00000043  8900              mov [bx+si],ax
00000045  8909              mov [bx+di],cx
00000047  8912              mov [bp+si],dx
00000049  891B              mov [bp+di],bx
0000004B  8924              mov [si],sp
0000004D  892D              mov [di],bp
0000004F  893F              mov [bx],di
00000051  894001            mov [bx+si+0x1],ax
00000054  8949FF            mov [bx+di-0x1],cx
00000057  895202            mov [bp+si+0x2],dx
0000005A  895BFE            mov [bp+di-0x2],bx
0000005D  896464            mov [si+0x64],sp
00000060  896D9C            mov [di-0x64],bp
00000063  897600            mov [bp+0x0],si
00000066  897601            mov [bp+0x1],si
00000069  897F01            mov [bx+0x1],di
0000006C  89800001          mov [bx+si+0x100],ax
00000070  898900FF          mov [bx+di-0x100],cx
00000074  89920002          mov [bp+si+0x200],dx
00000078  899B00FE          mov [bp+di-0x200],bx
0000007C  89A40064          mov [si+0x6400],sp
00000080  89AD009C          mov [di-0x6400],bp
00000084  89B60000          mov [bp+0x0],si
00000088  89B60001          mov [bp+0x100],si
0000008C  89BF0001          mov [bx+0x100],di
00000090  89C0              mov ax,ax
00000092  89C1              mov cx,ax
00000094  89C2              mov dx,ax
00000096  89C3              mov bx,ax
00000098  89C4              mov sp,ax
0000009A  89C5              mov bp,ax
0000009C  89C6              mov si,ax
0000009E  89C7              mov di,ax
000000A0  88C0              mov al,al
000000A2  88C1              mov cl,al
000000A4  88C2              mov dl,al
000000A6  88C3              mov bl,al
000000A8  88C4              mov ah,al
000000AA  88C5              mov ch,al
000000AC  88C6              mov dh,al
000000AE  88C7              mov bh,al
000000B0  C60012            mov byte [bx+si],0x12
000000B3  C606123456        mov byte [0x3412],0x56
000000B8  C6401234          mov byte [bx+si+0x12],0x34
000000BC  C680123456        mov byte [bx+si+0x3412],0x56
000000C1  C6C012            mov al,0x12
000000C4  C7001234          mov word [bx+si],0x3412
000000C8  C70612345678      mov word [0x3412],0x7856
000000CE  C740123456        mov word [bx+si+0x12],0x5634
000000D3  C78012345678      mov word [bx+si+0x3412],0x7856
000000D9  C7C01234          mov ax,0x3412
000000DD  A03412            mov al,[0x1234]
000000E0  A13412            mov ax,[0x1234]
000000E3  A23412            mov [0x1234],al
000000E6  A33412            mov [0x1234],ax
000000E9  8E00              mov es,[bx+si]
000000EB  8E4810            mov cs,[bx+si+0x10]
000000EE  8E9000F0          mov ss,[bx+si-0x1000]
000000F2  8ED8              mov ds,ax
000000F4  8C00              mov [bx+si],es
000000F6  8C4810            mov [bx+si+0x10],cs
000000F9  8C9000F0          mov [bx+si-0x1000],ss
000000FD  8CD8              mov ax,ds
000000FF  268800            mov [es:bx+si],al
00000102  2E88063412        mov [cs:0x1234],al
00000107  36894001          mov [ss:bx+si+0x1],ax
0000010B  3EC60212          mov byte [ds:bp+si],0x12
0000010F  268C9000F0        mov [es:bx+si-0x1000],ss
00000114  26FF30            push word [es:bx+si]
00000117  2E8F00            pop word [cs:bx+si]
0000011A  FF30              push word [bx+si]
0000011C  FF31              push word [bx+di]
0000011E  FF32              push word [bp+si]
00000120  FF33              push word [bp+di]
00000122  FF34              push word [si]
00000124  FF35              push word [di]
00000126  FF363412          push word [0x1234]
0000012A  FF37              push word [bx]
0000012C  FF7012            push word [bx+si+0x12]
0000012F  FF7112            push word [bx+di+0x12]
00000132  FF7212            push word [bp+si+0x12]
00000135  FF7312            push word [bp+di+0x12]
00000138  FF74FD            push word [si-0x3]
0000013B  FF75FD            push word [di-0x3]
0000013E  FF76FD            push word [bp-0x3]
00000141  FF77FD            push word [bx-0x3]
00000144  FFB03412          push word [bx+si+0x1234]
00000148  FFB13412          push word [bx+di+0x1234]
0000014C  FFB23412          push word [bp+si+0x1234]
00000150  FFB33412          push word [bp+di+0x1234]
00000154  FFB4FDFD          push word [si-0x203]
00000158  FFB5FDFD          push word [di-0x203]
0000015C  FFB6FDFD          push word [bp-0x203]
00000160  FFB7FDFD          push word [bx-0x203]
00000164  FFF0              push ax
00000166  FFF1              push cx
00000168  FFF2              push dx
0000016A  FFF3              push bx
0000016C  FFF4              push sp
0000016E  FFF5              push bp
00000170  FFF6              push si
00000172  FFF7              push di
00000174  50                push ax
00000175  51                push cx
00000176  52                push dx
00000177  53                push bx
00000178  54                push sp
00000179  55                push bp
0000017A  56                push si
0000017B  57                push di
0000017C  06                push es
0000017D  0E                push cs
0000017E  16                push ss
0000017F  1E                push ds
00000180  8F00              pop word [bx+si]
00000182  8F01              pop word [bx+di]
00000184  8F02              pop word [bp+si]
00000186  8F03              pop word [bp+di]
00000188  8F04              pop word [si]
0000018A  8F05              pop word [di]
0000018C  8F063412          pop word [0x1234]
00000190  8F07              pop word [bx]
00000192  8F4012            pop word [bx+si+0x12]
00000195  8F4112            pop word [bx+di+0x12]
00000198  8F4212            pop word [bp+si+0x12]
0000019B  8F4312            pop word [bp+di+0x12]
0000019E  8F44FD            pop word [si-0x3]
000001A1  8F45FD            pop word [di-0x3]
000001A4  8F46FD            pop word [bp-0x3]
000001A7  8F47FD            pop word [bx-0x3]
000001AA  8F803412          pop word [bx+si+0x1234]
000001AE  8F813412          pop word [bx+di+0x1234]
000001B2  8F823412          pop word [bp+si+0x1234]
000001B6  8F833412          pop word [bp+di+0x1234]
000001BA  8F84FDFD          pop word [si-0x203]
000001BE  8F85FDFD          pop word [di-0x203]
000001C2  8F86FDFD          pop word [bp-0x203]
000001C6  8F87FDFD          pop word [bx-0x203]
000001CA  8FC0              pop ax
000001CC  8FC1              pop cx
000001CE  8FC2              pop dx
000001D0  8FC3              pop bx
000001D2  8FC4              pop sp
000001D4  8FC5              pop bp
000001D6  8FC6              pop si
000001D8  8FC7              pop di
000001DA  58                pop ax
000001DB  59                pop cx
000001DC  5A                pop dx
000001DD  5B                pop bx
000001DE  5C                pop sp
000001DF  5D                pop bp
000001E0  5E                pop si
000001E1  5F                pop di
000001E2  07                pop es
000001E3  17                pop ss
000001E4  1F                pop ds
000001E5  8600              xchg al,[bx+si]
000001E7  8601              xchg al,[bx+di]
000001E9  8602              xchg al,[bp+si]
000001EB  8603              xchg al,[bp+di]
000001ED  8604              xchg al,[si]
000001EF  8605              xchg al,[di]
000001F1  86063412          xchg al,[0x1234]
000001F5  8607              xchg al,[bx]
000001F7  864812            xchg cl,[bx+si+0x12]
000001FA  864912            xchg cl,[bx+di+0x12]
000001FD  864A12            xchg cl,[bp+si+0x12]
00000200  864B12            xchg cl,[bp+di+0x12]
00000203  864CFD            xchg cl,[si-0x3]
00000206  864DFD            xchg cl,[di-0x3]
00000209  864EFD            xchg cl,[bp-0x3]
0000020C  864FFD            xchg cl,[bx-0x3]
0000020F  86903412          xchg dl,[bx+si+0x1234]
00000213  86913412          xchg dl,[bx+di+0x1234]
00000217  86923412          xchg dl,[bp+si+0x1234]
0000021B  86933412          xchg dl,[bp+di+0x1234]
0000021F  8694FDFD          xchg dl,[si-0x203]
00000223  8695FDFD          xchg dl,[di-0x203]
00000227  8696FDFD          xchg dl,[bp-0x203]
0000022B  8697FDFD          xchg dl,[bx-0x203]
0000022F  86D8              xchg bl,al
00000231  86D9              xchg bl,cl
00000233  86DA              xchg bl,dl
00000235  86DB              xchg bl,bl
00000237  86DC              xchg bl,ah
00000239  86DD              xchg bl,ch
0000023B  86DE              xchg bl,dh
0000023D  86DF              xchg bl,bh
0000023F  8700              xchg ax,[bx+si]
00000241  874812            xchg cx,[bx+si+0x12]
00000244  87903412          xchg dx,[bx+si+0x1234]
00000248  87D8              xchg bx,ax
0000024A  91                xchg ax,cx
0000024B  92                xchg ax,dx
0000024C  93                xchg ax,bx
0000024D  94                xchg ax,sp
0000024E  95                xchg ax,bp
0000024F  96                xchg ax,si
00000250  97                xchg ax,di
00000251  E412              in al,0x12
00000253  E534              in ax,0x34
00000255  EC                in al,dx
00000256  ED                in ax,dx
00000257  E6FF              out 0xff,al
00000259  E701              out 0x1,ax
0000025B  EE                out dx,al
0000025C  EF                out dx,ax
0000025D  D7                xlatb
0000025E  8D00              lea ax,[bx+si]
00000260  8D01              lea ax,[bx+di]
00000262  8D02              lea ax,[bp+si]
00000264  8D03              lea ax,[bp+di]
00000266  8D04              lea ax,[si]
00000268  8D05              lea ax,[di]
0000026A  8D063412          lea ax,[0x1234]
0000026E  8D07              lea ax,[bx]
00000270  8D4812            lea cx,[bx+si+0x12]
00000273  8D4912            lea cx,[bx+di+0x12]
00000276  8D4A12            lea cx,[bp+si+0x12]
00000279  8D4B12            lea cx,[bp+di+0x12]
0000027C  8D4CFD            lea cx,[si-0x3]
0000027F  8D4DFD            lea cx,[di-0x3]
00000282  8D4EFD            lea cx,[bp-0x3]
00000285  8D4FFD            lea cx,[bx-0x3]
00000288  8D903412          lea dx,[bx+si+0x1234]
0000028C  8D913412          lea dx,[bx+di+0x1234]
00000290  8D923412          lea dx,[bp+si+0x1234]
00000294  8D933412          lea dx,[bp+di+0x1234]
00000298  8D94FDFD          lea dx,[si-0x203]
0000029C  8D95FDFD          lea dx,[di-0x203]
000002A0  8D96FDFD          lea dx,[bp-0x203]
000002A4  8D97FDFD          lea dx,[bx-0x203]
000002A8  C500              lds ax,[bx+si]
000002AA  C501              lds ax,[bx+di]
000002AC  C502              lds ax,[bp+si]
000002AE  C503              lds ax,[bp+di]
000002B0  C504              lds ax,[si]
000002B2  C505              lds ax,[di]
000002B4  C5063412          lds ax,[0x1234]
000002B8  C507              lds ax,[bx]
000002BA  C54812            lds cx,[bx+si+0x12]
000002BD  C54912            lds cx,[bx+di+0x12]
000002C0  C54A12            lds cx,[bp+si+0x12]
000002C3  C54B12            lds cx,[bp+di+0x12]
000002C6  C54CFD            lds cx,[si-0x3]
000002C9  C54DFD            lds cx,[di-0x3]
000002CC  C54EFD            lds cx,[bp-0x3]
000002CF  C54FFD            lds cx,[bx-0x3]
000002D2  C5903412          lds dx,[bx+si+0x1234]
000002D6  C5913412          lds dx,[bx+di+0x1234]
000002DA  C5923412          lds dx,[bp+si+0x1234]
000002DE  C5933412          lds dx,[bp+di+0x1234]
000002E2  C594FDFD          lds dx,[si-0x203]
000002E6  C595FDFD          lds dx,[di-0x203]
000002EA  C596FDFD          lds dx,[bp-0x203]
000002EE  C597FDFD          lds dx,[bx-0x203]
000002F2  C400              les ax,[bx+si]
000002F4  C401              les ax,[bx+di]
000002F6  C402              les ax,[bp+si]
000002F8  C403              les ax,[bp+di]
000002FA  C404              les ax,[si]
000002FC  C405              les ax,[di]
000002FE  C4063412          les ax,[0x1234]
00000302  C407              les ax,[bx]
00000304  C44812            les cx,[bx+si+0x12]
00000307  C44912            les cx,[bx+di+0x12]
0000030A  C44A12            les cx,[bp+si+0x12]
0000030D  C44B12            les cx,[bp+di+0x12]
00000310  C44CFD            les cx,[si-0x3]
00000313  C44DFD            les cx,[di-0x3]
00000316  C44EFD            les cx,[bp-0x3]
00000319  C44FFD            les cx,[bx-0x3]
0000031C  C4903412          les dx,[bx+si+0x1234]
00000320  C4913412          les dx,[bx+di+0x1234]
00000324  C4923412          les dx,[bp+si+0x1234]
00000328  C4933412          les dx,[bp+di+0x1234]
0000032C  C494FDFD          les dx,[si-0x203]
00000330  C495FDFD          les dx,[di-0x203]
00000334  C496FDFD          les dx,[bp-0x203]
00000338  C497FDFD          les dx,[bx-0x203]
0000033C  9F                lahf
0000033D  9E                sahf
0000033E  9C                pushfw
0000033F  9D                popfw
00000340  0000              add [bx+si],al
00000342  0001              add [bx+di],al
00000344  0002              add [bp+si],al
00000346  0003              add [bp+di],al
00000348  0004              add [si],al
0000034A  0005              add [di],al
0000034C  00063412          add [0x1234],al
00000350  0007              add [bx],al
00000352  004812            add [bx+si+0x12],cl
00000355  004912            add [bx+di+0x12],cl
00000358  004A12            add [bp+si+0x12],cl
0000035B  004B12            add [bp+di+0x12],cl
0000035E  004CFD            add [si-0x3],cl
00000361  004DFD            add [di-0x3],cl
00000364  004EFD            add [bp-0x3],cl
00000367  004FFD            add [bx-0x3],cl
0000036A  00903412          add [bx+si+0x1234],dl
0000036E  00913412          add [bx+di+0x1234],dl
00000372  00923412          add [bp+si+0x1234],dl
00000376  00933412          add [bp+di+0x1234],dl
0000037A  0094FDFD          add [si-0x203],dl
0000037E  0095FDFD          add [di-0x203],dl
00000382  0096FDFD          add [bp-0x203],dl
00000386  0097FDFD          add [bx-0x203],dl
0000038A  00C0              add al,al
0000038C  00C1              add cl,al
0000038E  00C2              add dl,al
00000390  00C3              add bl,al
00000392  00C4              add ah,al
00000394  00C5              add ch,al
00000396  00C6              add dh,al
00000398  00C7              add bh,al
0000039A  0100              add [bx+si],ax
0000039C  0101              add [bx+di],ax
0000039E  0102              add [bp+si],ax
000003A0  0103              add [bp+di],ax
000003A2  0104              add [si],ax
000003A4  0105              add [di],ax
000003A6  01063412          add [0x1234],ax
000003AA  0107              add [bx],ax
000003AC  014812            add [bx+si+0x12],cx
000003AF  014912            add [bx+di+0x12],cx
000003B2  014A12            add [bp+si+0x12],cx
000003B5  014B12            add [bp+di+0x12],cx
000003B8  014CFD            add [si-0x3],cx
000003BB  014DFD            add [di-0x3],cx
000003BE  014EFD            add [bp-0x3],cx
000003C1  014FFD            add [bx-0x3],cx
000003C4  01903412          add [bx+si+0x1234],dx
000003C8  01913412          add [bx+di+0x1234],dx
000003CC  01923412          add [bp+si+0x1234],dx
000003D0  01933412          add [bp+di+0x1234],dx
000003D4  0194FDFD          add [si-0x203],dx
000003D8  0195FDFD          add [di-0x203],dx
000003DC  0196FDFD          add [bp-0x203],dx
000003E0  0197FDFD          add [bx-0x203],dx
000003E4  01C0              add ax,ax
000003E6  01C1              add cx,ax
000003E8  01C2              add dx,ax
000003EA  01C3              add bx,ax
000003EC  01C4              add sp,ax
000003EE  01C5              add bp,ax
000003F0  01C6              add si,ax
000003F2  01C7              add di,ax
000003F4  0200              add al,[bx+si]
000003F6  0201              add al,[bx+di]
000003F8  0202              add al,[bp+si]
000003FA  0203              add al,[bp+di]
000003FC  0204              add al,[si]
000003FE  0205              add al,[di]
00000400  02063412          add al,[0x1234]
00000404  0207              add al,[bx]
00000406  024812            add cl,[bx+si+0x12]
00000409  024912            add cl,[bx+di+0x12]
0000040C  024A12            add cl,[bp+si+0x12]
0000040F  024B12            add cl,[bp+di+0x12]
00000412  024CFD            add cl,[si-0x3]
00000415  024DFD            add cl,[di-0x3]
00000418  024EFD            add cl,[bp-0x3]
0000041B  024FFD            add cl,[bx-0x3]
0000041E  02903412          add dl,[bx+si+0x1234]
00000422  02913412          add dl,[bx+di+0x1234]
00000426  02923412          add dl,[bp+si+0x1234]
0000042A  02933412          add dl,[bp+di+0x1234]
0000042E  0294FDFD          add dl,[si-0x203]
00000432  0295FDFD          add dl,[di-0x203]
00000436  0296FDFD          add dl,[bp-0x203]
0000043A  0297FDFD          add dl,[bx-0x203]
0000043E  02C0              add al,al
00000440  02C1              add al,cl
00000442  02C2              add al,dl
00000444  02C3              add al,bl
00000446  02C4              add al,ah
00000448  02C5              add al,ch
0000044A  02C6              add al,dh
0000044C  02C7              add al,bh
0000044E  0300              add ax,[bx+si]
00000450  0301              add ax,[bx+di]
00000452  0302              add ax,[bp+si]
00000454  0303              add ax,[bp+di]
00000456  0304              add ax,[si]
00000458  0305              add ax,[di]
0000045A  03063412          add ax,[0x1234]
0000045E  0307              add ax,[bx]
00000460  034812            add cx,[bx+si+0x12]
00000463  034912            add cx,[bx+di+0x12]
00000466  034A12            add cx,[bp+si+0x12]
00000469  034B12            add cx,[bp+di+0x12]
0000046C  034CFD            add cx,[si-0x3]
0000046F  034DFD            add cx,[di-0x3]
00000472  034EFD            add cx,[bp-0x3]
00000475  034FFD            add cx,[bx-0x3]
00000478  03903412          add dx,[bx+si+0x1234]
0000047C  03913412          add dx,[bx+di+0x1234]
00000480  03923412          add dx,[bp+si+0x1234]
00000484  03933412          add dx,[bp+di+0x1234]
00000488  0394FDFD          add dx,[si-0x203]
0000048C  0395FDFD          add dx,[di-0x203]
00000490  0396FDFD          add dx,[bp-0x203]
00000494  0397FDFD          add dx,[bx-0x203]
00000498  03C0              add ax,ax
0000049A  03C1              add ax,cx
0000049C  03C2              add ax,dx
0000049E  03C3              add ax,bx
000004A0  03C4              add ax,sp
000004A2  03C5              add ax,bp
000004A4  03C6              add ax,si
000004A6  03C7              add ax,di
000004A8  800000            add byte [bx+si],0x0
000004AB  8100FFFF          add word [bx+si],0xffff
000004AF  830000            add word [bx+si],byte +0x0
000004B2  8300FF            add word [bx+si],byte -0x1
000004B5  0400              add al,0x0
000004B7  0534FF            add ax,0xff34
000004BA  1000              adc [bx+si],al
000004BC  1001              adc [bx+di],al
000004BE  1002              adc [bp+si],al
000004C0  1003              adc [bp+di],al
000004C2  1004              adc [si],al
000004C4  1005              adc [di],al
000004C6  10063412          adc [0x1234],al
000004CA  1007              adc [bx],al
000004CC  104812            adc [bx+si+0x12],cl
000004CF  104912            adc [bx+di+0x12],cl
000004D2  104A12            adc [bp+si+0x12],cl
000004D5  104B12            adc [bp+di+0x12],cl
000004D8  104CFD            adc [si-0x3],cl
000004DB  104DFD            adc [di-0x3],cl
000004DE  104EFD            adc [bp-0x3],cl
000004E1  104FFD            adc [bx-0x3],cl
000004E4  10903412          adc [bx+si+0x1234],dl
000004E8  10913412          adc [bx+di+0x1234],dl
000004EC  10923412          adc [bp+si+0x1234],dl
000004F0  10933412          adc [bp+di+0x1234],dl
000004F4  1094FDFD          adc [si-0x203],dl
000004F8  1095FDFD          adc [di-0x203],dl
000004FC  1096FDFD          adc [bp-0x203],dl
00000500  1097FDFD          adc [bx-0x203],dl
00000504  10C0              adc al,al
00000506  10C1              adc cl,al
00000508  10C2              adc dl,al
0000050A  10C3              adc bl,al
0000050C  10C4              adc ah,al
0000050E  10C5              adc ch,al
00000510  10C6              adc dh,al
00000512  10C7              adc bh,al
00000514  1100              adc [bx+si],ax
00000516  1101              adc [bx+di],ax
00000518  1102              adc [bp+si],ax
0000051A  1103              adc [bp+di],ax
0000051C  1104              adc [si],ax
0000051E  1105              adc [di],ax
00000520  11063412          adc [0x1234],ax
00000524  1107              adc [bx],ax
00000526  114812            adc [bx+si+0x12],cx
00000529  114912            adc [bx+di+0x12],cx
0000052C  114A12            adc [bp+si+0x12],cx
0000052F  114B12            adc [bp+di+0x12],cx
00000532  114CFD            adc [si-0x3],cx
00000535  114DFD            adc [di-0x3],cx
00000538  114EFD            adc [bp-0x3],cx
0000053B  114FFD            adc [bx-0x3],cx
0000053E  11903412          adc [bx+si+0x1234],dx
00000542  11913412          adc [bx+di+0x1234],dx
00000546  11923412          adc [bp+si+0x1234],dx
0000054A  11933412          adc [bp+di+0x1234],dx
0000054E  1194FDFD          adc [si-0x203],dx
00000552  1195FDFD          adc [di-0x203],dx
00000556  1196FDFD          adc [bp-0x203],dx
0000055A  1197FDFD          adc [bx-0x203],dx
0000055E  11C0              adc ax,ax
00000560  11C1              adc cx,ax
00000562  11C2              adc dx,ax
00000564  11C3              adc bx,ax
00000566  11C4              adc sp,ax
00000568  11C5              adc bp,ax
0000056A  11C6              adc si,ax
0000056C  11C7              adc di,ax
0000056E  1200              adc al,[bx+si]
00000570  1201              adc al,[bx+di]
00000572  1202              adc al,[bp+si]
00000574  1203              adc al,[bp+di]
00000576  1204              adc al,[si]
00000578  1205              adc al,[di]
0000057A  12063412          adc al,[0x1234]
0000057E  1207              adc al,[bx]
00000580  124812            adc cl,[bx+si+0x12]
00000583  124912            adc cl,[bx+di+0x12]
00000586  124A12            adc cl,[bp+si+0x12]
00000589  124B12            adc cl,[bp+di+0x12]
0000058C  124CFD            adc cl,[si-0x3]
0000058F  124DFD            adc cl,[di-0x3]
00000592  124EFD            adc cl,[bp-0x3]
00000595  124FFD            adc cl,[bx-0x3]
00000598  12903412          adc dl,[bx+si+0x1234]
0000059C  12913412          adc dl,[bx+di+0x1234]
000005A0  12923412          adc dl,[bp+si+0x1234]
000005A4  12933412          adc dl,[bp+di+0x1234]
000005A8  1294FDFD          adc dl,[si-0x203]
000005AC  1295FDFD          adc dl,[di-0x203]
000005B0  1296FDFD          adc dl,[bp-0x203]
000005B4  1297FDFD          adc dl,[bx-0x203]
000005B8  12C0              adc al,al
000005BA  12C1              adc al,cl
000005BC  12C2              adc al,dl
000005BE  12C3              adc al,bl
000005C0  12C4              adc al,ah
000005C2  12C5              adc al,ch
000005C4  12C6              adc al,dh
000005C6  12C7              adc al,bh
000005C8  1300              adc ax,[bx+si]
000005CA  1301              adc ax,[bx+di]
000005CC  1302              adc ax,[bp+si]
000005CE  1303              adc ax,[bp+di]
000005D0  1304              adc ax,[si]
000005D2  1305              adc ax,[di]
000005D4  13063412          adc ax,[0x1234]
000005D8  1307              adc ax,[bx]
000005DA  134812            adc cx,[bx+si+0x12]
000005DD  134912            adc cx,[bx+di+0x12]
000005E0  134A12            adc cx,[bp+si+0x12]
000005E3  134B12            adc cx,[bp+di+0x12]
000005E6  134CFD            adc cx,[si-0x3]
000005E9  134DFD            adc cx,[di-0x3]
000005EC  134EFD            adc cx,[bp-0x3]
000005EF  134FFD            adc cx,[bx-0x3]
000005F2  13903412          adc dx,[bx+si+0x1234]
000005F6  13913412          adc dx,[bx+di+0x1234]
000005FA  13923412          adc dx,[bp+si+0x1234]
000005FE  13933412          adc dx,[bp+di+0x1234]
00000602  1394FDFD          adc dx,[si-0x203]
00000606  1395FDFD          adc dx,[di-0x203]
0000060A  1396FDFD          adc dx,[bp-0x203]
0000060E  1397FDFD          adc dx,[bx-0x203]
00000612  13C0              adc ax,ax
00000614  13C1              adc ax,cx
00000616  13C2              adc ax,dx
00000618  13C3              adc ax,bx
0000061A  13C4              adc ax,sp
0000061C  13C5              adc ax,bp
0000061E  13C6              adc ax,si
00000620  13C7              adc ax,di
00000622  801000            adc byte [bx+si],0x0
00000625  8110FFFF          adc word [bx+si],0xffff
00000629  831000            adc word [bx+si],byte +0x0
0000062C  8310FF            adc word [bx+si],byte -0x1
0000062F  1400              adc al,0x0
00000631  1534FF            adc ax,0xff34
00000634  8006123456        add byte [0x3412],0x56
00000639  C4                db 0xc4
0000063A  FF                db 0xff
0000063B  FE00              inc byte [bx+si]
0000063D  FF00              inc word [bx+si]
0000063F  FE873456          inc byte [bx+0x5634]
00000643  FEC7              inc bh
00000645  40                inc ax
00000646  41                inc cx
00000647  42                inc dx
00000648  37                aaa
00000649  27                daa
0000064A  2800              sub [bx+si],al
0000064C  2900              sub [bx+si],ax
0000064E  2A00              sub al,[bx+si]
00000650  8028FF            sub byte [bx+si],0xff
00000653  2CFF              sub al,0xff
00000655  2DFFFF            sub ax,0xffff
00000658  18FE              sbb dh,bh
0000065A  19FE              sbb si,di
0000065C  1AFE              sbb bh,dh
0000065E  80583456          sbb byte [bx+si+0x34],0x56
00000662  1C34              sbb al,0x34
00000664  1DFFFF            sbb ax,0xffff
00000667  FEC8              dec al
00000669  FF0EAAAA          dec word [0xaaaa]
0000066D  49                dec cx
0000066E  4D                dec bp
0000066F  F618              neg byte [bx+si]
00000671  F7D8              neg ax
00000673  3800              cmp [bx+si],al
00000675  3B00              cmp ax,[bx+si]
00000677  8038AA            cmp byte [bx+si],0xaa
0000067A  8338AA            cmp word [bx+si],byte -0x56
0000067D  3CBB              cmp al,0xbb
0000067F  3D4567            cmp ax,0x6745
00000682  3F                aas
00000683  2F                das
00000684  F620              mul byte [bx+si]
00000686  F720              mul word [bx+si]
00000688  F628              imul byte [bx+si]
0000068A  F728              imul word [bx+si]
0000068C  D40A              aam
0000068E  F630              div byte [bx+si]
00000690  F730              div word [bx+si]
00000692  F738              idiv word [bx+si]
00000694  F738              idiv word [bx+si]
00000696  D50A              aad
00000698  98                cbw
00000699  99                cwd
0000069A  F610              not byte [bx+si]
0000069C  D020              shl byte [bx+si],1
0000069E  D120              shl word [bx+si],1
000006A0  D220              shl byte [bx+si],cl
000006A2  D320              shl word [bx+si],cl
000006A4  D060AA            shl byte [bx+si-0x56],1
000006A7  D026AAAA          shl byte [0xaaaa],1
000006AB  D028              shr byte [bx+si],1
000006AD  D328              shr word [bx+si],cl
000006AF  D038              sar byte [bx+si],1
000006B1  D338              sar word [bx+si],cl
000006B3  D000              rol byte [bx+si],1
000006B5  D300              rol word [bx+si],cl
000006B7  D008              ror byte [bx+si],1
000006B9  D308              ror word [bx+si],cl
000006BB  D010              rcl byte [bx+si],1
000006BD  D310              rcl word [bx+si],cl
000006BF  D018              rcr byte [bx+si],1
000006C1  D318              rcr word [bx+si],cl
000006C3  2000              and [bx+si],al
000006C5  2244AA            and al,[si-0x56]
000006C8  81203412          and word [bx+si],0x1234
000006CC  25FFFF            and ax,0xffff
000006CF  8500              test [bx+si],ax
000006D1  F60034            test byte [bx+si],0x34
000006D4  A8DD              test al,0xdd
000006D6  0B00              or ax,[bx+si]
000006D8  81083456          or word [bx+si],0x5634
000006DC  0CDD              or al,0xdd
000006DE  3000              xor [bx+si],al
000006E0  81F7AAAA          xor di,0xaaaa
000006E4  34AA              xor al,0xaa
000006E6  A4                movsb
000006E7  A5                movsw
000006E8  A6                cmpsb
000006E9  A7                cmpsw
000006EA  AE                scasb
000006EB  AF                scasw
000006EC  AC                lodsb
000006ED  AD                lodsw
000006EE  AA                stosb
000006EF  AB                stosw
000006F0  F2A4              repne movsb
000006F2  F3AB              rep stosw
000006F4  E80000            call word 0x6f7
000006F7  FF10              call word [bx+si]
000006F9  FF11              call word [bx+di]
000006FB  FF12              call word [bp+si]
000006FD  FF13              call word [bp+di]
000006FF  FF14              call word [si]
00000701  FF15              call word [di]
00000703  FF161234          call word [0x3412]
00000707  FF17              call word [bx]
00000709  FF5034            call word [bx+si+0x34]
0000070C  FF51FF            call word [bx+di-0x1]
0000070F  FF5234            call word [bp+si+0x34]
00000712  FF53FF            call word [bp+di-0x1]
00000715  FF5434            call word [si+0x34]
00000718  FF55FF            call word [di-0x1]
0000071B  FF5634            call word [bp+0x34]
0000071E  FF57FF            call word [bx-0x1]
00000721  FFD0              call ax
00000723  FFD1              call cx
00000725  FFD2              call dx
00000727  FFD3              call bx
00000729  FFD4              call sp
0000072B  FFD5              call bp
0000072D  FFD6              call si
0000072F  FFD7              call di
00000731  FF18              call word far [bx+si]
00000733  FF19              call word far [bx+di]
00000735  FF1A              call word far [bp+si]
00000737  FF1B              call word far [bp+di]
00000739  FF1C              call word far [si]
0000073B  FF1D              call word far [di]
0000073D  FF1E3412          call word far [0x1234]
00000741  FF1F              call word far [bx]
00000743  FF58FF            call word far [bx+si-0x1]
00000746  FF983456          call word far [bx+si+0x5634]
0000074A  9A12345678        call word 0x7856:0x3412
0000074F  E90000            jmp word 0x752
00000752  EB00              jmp short 0x754
00000754  FF20              jmp word [bx+si]
00000756  FF60FF            jmp word [bx+si-0x1]
00000759  FFA0FFFF          jmp word [bx+si-0x1]
0000075D  FFE0              jmp ax
0000075F  EA01234567        jmp word 0x6745:0x2301
00000764  FF28              jmp word far [bx+si]
00000766  FFA8AAAA          jmp word far [bx+si-0x5556]
0000076A  C3                ret
0000076B  C21234            ret 0x3412
0000076E  CB                retf
0000076F  CA0000            retf 0x0
00000772  7400              jz 0x774
00000774  7C00              jl 0x776
00000776  7E00              jng 0x778
00000778  7200              jc 0x77a
0000077A  7600              jna 0x77c
0000077C  7A00              jpe 0x77e
0000077E  7000              jo 0x780
00000780  7800              js 0x782
00000782  7500              jnz 0x784
00000784  7D00              jnl 0x786
00000786  7F00              jg 0x788
00000788  7300              jnc 0x78a
0000078A  7700              ja 0x78c
0000078C  7B00              jpo 0x78e
0000078E  7100              jno 0x790
00000790  7900              jns 0x792
00000792  E200              loop 0x794
00000794  E100              loope 0x796
00000796  E000              loopne 0x798
00000798  E300              jcxz 0x79a
0000079A  CDFF              int 0xff
0000079C  CC                int3
0000079D  CE                into
0000079E  CF                iretw
0000079F  F8                clc
000007A0  F5                cmc
000007A1  F9                stc
000007A2  FC                cld
000007A3  FD                std
000007A4  FA                cli
000007A5  FB                sti
000007A6  F4                hlt
000007A7  FE06FFFF          inc byte [0xffff]
000007AB  268CD8            es mov ax,ds
000007AE  26F4              es hlt
000007B0  260400            es add al,0x0
000007B3  26FE06FFFF        inc byte [es:0xffff]
000007B8  2640              es inc ax
000007BA  D1FF              sar di,1
000007BC  90                nop
000007BD  F00000            lock add [bx+si],al
000007C0  9B0000            wait add [bx+si],al
000007C3  EBFD              jmp short 0x7c2
000007C5  E9FFFD            jmp word 0x5c7
000007C8  EBF0              jmp short 0x7ba
