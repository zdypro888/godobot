![](media/image1.jpeg){width="1.4749989063867017in"
height="0.5583333333333333in"}

> [**接口说明**]{.underline}
>
> **Dobot Magician API**
>
> **接口说明**
>
> 文档版本： V1.2.3
>
> 发布日期：2019-07-19
>
> 深圳市越疆科技有限公司
>
> **版权所有** **© 越疆科技有限公司2017。** **保留一切权利。**
>
> 非经本公司书面许可，任何单位和个人不得擅自摘抄、复制本文档内容的部分或全
> 部，并不得以任何形式传播。
>
> **免责申明**
>
> 在法律允许的最大范围内，本手册所描述的产品（含其硬件、软件、固件等）均"按
> 照现状"提供， 可能存在瑕疵、错误或故障，
> 越疆不提供任何形式的明示或默示保证，
> 包括但不限于适销性、质量满意度、适合特定目的、不侵犯第三方权利等保证；
> 亦不对
> 使用本手册或使用本公司产品导致的任何特殊、附带、偶然或间接的损害进行赔偿。
>
> 在使用本产品前详细阅读本使用手册及网上发布的相关技术文档并了解相关信息，
> 确保在充分了解机器人及其相关知识的前提下使用机械臂。越疆建议您在专业人员的指
> 导下使用本手册。该手册所包含的所有安全方面的信息都不得视为Dobot的保证，即便
> 遵循本手册及相关说明， 使用过程中造成的危害或损失依然有可能发生。
>
> 本产品的使用者有责任确保遵循相关国家的切实可行的法律法规，确保在越疆机械
> 臂的使用中不存在任何重大危险。
>
> 越疆科技有限公司
>
> 地址：深圳市南山区同富裕工业城三栋三楼
>
> 网址：[[http://cn.dobot.cc/]{.underline}](http://cn.dobot.cc/)
>
> **前** **言**
>
> **目的**
>
> 本文档旨在对 Dobot API 接口进行详细说明，并给出基于 Dobot API
> 接口开发应用程序 的一般流程。
>
> **读者对象**
>
> 本手册适用于：
>
> . 客户工程师
>
> . 安装调测工程师
>
> . 技术支持工程师
>
> **修订记录**

<table>
<colgroup>
<col style="width: 50%" />
<col style="width: 49%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>时间</p>
</blockquote></td>
<td><blockquote>
<p>修订记录</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>2019/07/19</p>
</blockquote></td>
<td><table>
<colgroup>
<col style="width: 7%" />
<col style="width: 92%" />
</colgroup>
<tbody>
<tr class="odd">
<td><p>.</p>
<p>.</p>
<p>.</p></td>
<td><blockquote>
<p>增加AIP返回结果说明</p>
<p>增加速度、速度比例之间关系说明</p>
<p>调整目录结构</p>
</blockquote></td>
</tr>
</tbody>
</table></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>2019/05/05</p>
</blockquote></td>
<td><blockquote>
<p>给SetJOGCmd接口函数添加说明</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>2018/11/06</p>
</blockquote></td>
<td><blockquote>
<p>修正了一些API的错误</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>2018/03/21</p>
</blockquote></td>
<td><blockquote>
<p>优化所有API</p>
</blockquote></td>
</tr>
</tbody>
</table>

> **符号约定**
>
> 在本手册中可能出现下列标志， 它们所代表的含义如下。

<table>
<colgroup>
<col style="width: 6%" />
<col style="width: 43%" />
<col style="width: 49%" />
</colgroup>
<tbody>
<tr class="odd">
<td colspan="2"><blockquote>
<p>符号</p>
</blockquote></td>
<td><blockquote>
<p>说明</p>
</blockquote></td>
</tr>
<tr class="even">
<td><img src="media/image4.png"
style="width:0.27726in;height:0.23938in" /></td>
<td><blockquote>
<p>危险</p>
</blockquote></td>
<td><blockquote>
<p>表示有高度潜在危险， 如果不能避免， 会导致 人员死亡或严重伤害</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><img src="media/image4.png"
style="width:0.27726in;height:0.23938in" /></td>
<td><blockquote>
<p>警告</p>
</blockquote></td>
<td><blockquote>
<p>表示有中度或低度潜在危害，如果不能避免，
可能导致人员轻微伤害、机械臂毁坏等情况</p>
</blockquote></td>
</tr>
<tr class="even">
<td><img src="media/image4.png"
style="width:0.27726in;height:0.23938in" /></td>
<td><blockquote>
<p>注意</p>
</blockquote></td>
<td><blockquote>
<p>表示有潜在风险， 如果忽视这些文本， 可能导
致机械臂损坏、数据丢失或不可预知的结果</p>
</blockquote></td>
</tr>
<tr class="odd">
<td colspan="2"><blockquote>
<p>说明</p>
</blockquote></td>
<td><blockquote>
<p>表示是正文的附加信息， 是对正文的强调和补 充</p>
</blockquote></td>
</tr>
</tbody>
</table>

> **目** **录**
>
> [![](media/image6.png){width="0.11499890638670167in"
> height="0.12333333333333334in"} **Dobot 指令简介** **1**](#_bookmark1)
>
> [![](media/image7.png){width="0.12166666666666667in"
> height="0.12333333333333334in"} **指令超时** **2**](#_bookmark2)
>
> [![](media/image8.png){width="0.17166666666666666in"
> height="0.10999890638670166in"} 设置指令超时时间 2](#_bookmark3)
>
> [![](media/image9.png){width="0.123332239720035in"
> height="0.12333333333333334in"} **连接/断开** **3**](#_bookmark4)
>
> [![](media/image10.png){width="0.16833333333333333in"
> height="0.10999890638670166in"} 搜索 Dobot 3](#_bookmark5)
>
> [![](media/image11.png){width="0.17999890638670166in"
> height="0.10999890638670166in"} 连接 Dobot 3](#_bookmark6)
>
> [![](media/image12.png){width="0.173332239720035in"
> height="0.10999890638670166in"} 断开 Dobot 4](#_bookmark7)
>
> [![](media/image13.png){width="0.17999890638670166in"
> height="0.10999890638670166in"} 示例： 连接示例 4](#_bookmark8)
>
> [![](media/image14.png){width="0.12166666666666667in"
> height="0.12333333333333334in"} **指令队列控制** **6**](#_bookmark9)
>
> [![](media/image15.png){width="0.17166666666666666in"
> height="0.10999890638670166in"} 执行队列中的指令 6](#_bookmark10)
>
> [![](media/image16.png){width="0.18333223972003498in"
> height="0.10999890638670166in"} 停止执行队列中的指令 6](#_bookmark11)
>
> [![](media/image17.png){width="0.17666666666666667in"
> height="0.10999890638670166in"} 强制停止执行队列中的指令
> 6](#_bookmark12)
>
> [![](media/image18.png){width="0.18333223972003498in"
> height="0.10999890638670166in"} 示例： 同步处理 PTP 指令下发和队列控制
> 7](#_bookmark13)
>
> [![](media/image19.png){width="0.17999890638670166in"
> height="0.10999890638670166in"} 示例： 异步处理 PTP 指令下发和队列控制
> 7](#_bookmark14)
>
> [![](media/image20.png){width="0.18333223972003498in"
> height="0.10999890638670166in"} 下载指令 8](#_bookmark15)
>
> [![](media/image21.png){width="0.18333223972003498in" height="0.11in"}
> 停止下载指令 9](#_bookmark16)
>
> [![](media/image22.png){width="0.17999890638670166in" height="0.11in"}
> 示例： 下载 PTP 指令 9](#_bookmark17)
>
> [![](media/image23.png){width="0.18333223972003498in" height="0.11in"}
> 清空指令队列 10](#_bookmark18)
>
> [![](media/image24.png){width="0.25666557305336835in" height="0.11in"}
> 获取指令索引 10](#_bookmark19)
>
> [![](media/image25.png){width="0.23999890638670165in" height="0.11in"}
> 示例： 获取指令索引实现运动同步 10](#_bookmark20)
>
> [![](media/image26.png){width="0.11999890638670166in"
> height="0.12166666666666667in"} **设备信息** **12**](#_bookmark21)
>
> [![](media/image27.png){width="0.16666666666666666in" height="0.11in"}
> 设置设备序列号 12](#_bookmark22)
>
> [![](media/image28.png){width="0.178332239720035in" height="0.11in"}
> 获取设备序列号 12](#_bookmark23)
>
> [![](media/image29.png){width="0.17166666666666666in" height="0.11in"}
> 设置设备名称 12](#_bookmark24)
>
> [![](media/image30.png){width="0.178332239720035in" height="0.11in"}
> 获取设备名称 12](#_bookmark25)
>
> [![](media/image31.png){width="0.17499890638670165in"
> height="0.108332239720035in"} 获取设备版本号 13](#_bookmark26)
>
> [![](media/image32.png){width="0.178332239720035in" height="0.11in"}
> 设置滑轨状态 13](#_bookmark27)
>
> [![](media/image33.png){width="0.178332239720035in"
> height="0.108332239720035in"} 获取滑轨状态 13](#_bookmark28)
>
> [![](media/image34.png){width="0.17499890638670165in" height="0.11in"}
> 获取设备时钟 14](#_bookmark29)
>
> [![](media/image35.png){width="0.11999890638670166in"
> height="0.12333114610673666in"} **实时位姿** **15**](#_bookmark30)
>
> [![](media/image36.png){width="0.16833333333333333in" height="0.11in"}
> 获取机械臂实时位姿 15](#_bookmark31)
>
> [![](media/image37.png){width="0.17999890638670166in"
> height="0.10999890638670166in"} 获取滑轨实时位置 15](#_bookmark32)
>
> [![](media/image38.png){width="0.173332239720035in"
> height="0.10999890638670166in"} 重设机器人实时位姿 15](#_bookmark33)
>
> [![](media/image39.png){width="0.11999890638670166in"
> height="0.12166447944006999in"} **报警功能** **17**](#_bookmark34)
>
> [![](media/image40.png){width="0.16833333333333333in"
> height="0.10999890638670166in"} 获取系统报警状态 17](#_bookmark35)
>
> [![](media/image41.png){width="0.17999890638670166in"
> height="0.10999890638670166in"} 清除系统所有报警 17](#_bookmark36)
>
> [![](media/image42.png){width="0.11999890638670166in"
> height="0.12333114610673666in"} **回零功能** **18**](#_bookmark37)
>
> [![](media/image43.png){width="0.16499890638670167in"
> height="0.10999890638670166in"} 设置回零位置 18](#_bookmark38)
>
> [![](media/image44.png){width="0.17666666666666667in"
> height="0.10999890638670166in"} 获取回零位置 18](#_bookmark39)
>
> [![](media/image45.png){width="0.16999890638670168in"
> height="0.10999890638670166in"} 执行回零功能 19](#_bookmark40)
>
> [![](media/image46.png){width="0.17666666666666667in"
> height="0.10999890638670166in"} 执行自动调平功能 19](#_bookmark41)
>
> [![](media/image47.png){width="0.17333333333333334in"
> height="0.10999890638670166in"} 获取自动调平结果 20](#_bookmark42)
>
> [![](media/image48.png){width="0.11999890638670166in"
> height="0.12333333333333334in"} **HHT 功能 21**](#_bookmark43)
>
> [![](media/image49.png){width="0.16833333333333333in"
> height="0.10999890638670166in"} 设置触发模式 21](#_bookmark44)
>
> [![](media/image50.png){width="0.17999890638670166in"
> height="0.10999890638670166in"} 获取触发模式 21](#_bookmark45)
>
> [![](media/image51.png){width="0.173332239720035in"
> height="0.10999890638670166in"} 设置手持示教使能状态 21](#_bookmark46)
>
> [![](media/image52.png){width="0.17999890638670166in"
> height="0.10999890638670166in"} 获取手持示教功能使能状态
> 22](#_bookmark47)
>
> [![](media/image53.png){width="0.17666666666666667in"
> height="0.10999890638670166in"} 获取手持示教触发信号 22](#_bookmark48)
>
> [![](media/image54.png){width="0.17999890638670166in"
> height="0.10999890638670166in"} 示例： 手持示教 22](#_bookmark49)
>
> [![](media/image55.png){width="0.198332239720035in"
> height="0.12333333333333334in"} **末端执行器** **24**](#_bookmark50)
>
> [![](media/image56.png){width="0.23in" height="0.10999890638670166in"}
> 设置末端坐标偏移量 24](#_bookmark51)
>
> [![](media/image57.png){width="0.24166666666666667in"
> height="0.10999890638670166in"} 获取末端坐标偏移量 24](#_bookmark52)
>
> [![](media/image58.png){width="0.235in"
> height="0.10999890638670166in"} 设置激光状态 25](#_bookmark53)
>
> [![](media/image59.png){width="0.24166666666666667in"
> height="0.10999890638670166in"} 获取激光状态 25](#_bookmark54)
>
> [![](media/image60.png){width="0.238332239720035in"
> height="0.10999890638670166in"} 设置气泵状态 25](#_bookmark55)
>
> [![](media/image61.png){width="0.24166666666666667in"
> height="0.10999890638670166in"} 获取气泵状态 26](#_bookmark56)
>
> [![](media/image62.png){width="0.24166666666666667in"
> height="0.10999890638670166in"} 设置夹爪状态 26](#_bookmark57)
>
> [![](media/image63.png){width="0.238332239720035in"
> height="0.10999890638670166in"} 获取夹爪状态 26](#_bookmark58)
>
> [![](media/image64.png){width="0.198332239720035in"
> height="0.12333114610673666in"} **JOG 功能** **28**](#_bookmark59)
>
> [![](media/image65.png){width="0.225in" height="0.11in"}
> 设置点动时各关节坐标轴的动速度和加速度 28](#_bookmark60)
>
> [![](media/image66.png){width="0.23666666666666666in" height="0.11in"}
> 获取点动时各关节坐标轴的动速度和加速度 28](#_bookmark61)
>
> [![](media/image67.png){width="0.23in" height="0.11in"}
> 设置点动时笛卡尔坐标轴的速度和加速度 29](#_bookmark62)
>
> [![](media/image68.png){width="0.23666666666666666in" height="0.11in"}
> 获取点动时笛卡尔坐标轴的速度和加速度 29](#_bookmark63)
>
> [![](media/image69.png){width="0.233332239720035in" height="0.11in"}
> 设置点动时滑轨速度和加速度 29](#_bookmark64)
>
> [![](media/image70.png){width="0.23666666666666666in" height="0.11in"}
> 获取点动时滑轨速度和加速度 30](#_bookmark65)
>
> [![](media/image71.png){width="0.23666666666666666in" height="0.11in"}
> 设置点动速度百分比和加速度百分比 30](#_bookmark66)
>
> [![](media/image72.png){width="0.233332239720035in" height="0.11in"}
> 获取点动速度百分比和加速度百分比 31](#_bookmark67)
>
> [![](media/image73.png){width="0.23666666666666666in" height="0.11in"}
> 执行点动指令 31](#_bookmark68)
>
> [![](media/image74.png){width="0.198332239720035in"
> height="0.12333333333333334in"} **PTP 功能 33**](#_bookmark69)
>
> [![](media/image75.png){width="0.23in" height="0.11in"} 设置 PTP
> 模式下各关节坐标轴的速度和加速度 33](#_bookmark70)
>
> [![](media/image76.png){width="0.24166666666666667in" height="0.11in"}
> 获取 PTP 模式下各关节坐标轴的速度和加速度 34](#_bookmark71)
>
> [![](media/image77.png){width="0.235in" height="0.11in"} 设置 PTP
> 模式下各笛卡尔坐标轴的速度和加速度 34](#_bookmark72)
>
> [![](media/image78.png){width="0.24166666666666667in" height="0.11in"}
> 获取 PTP 模式下各笛卡尔坐标轴的速度和加速度 35](#_bookmark73)
>
> [![](media/image79.png){width="0.238332239720035in" height="0.11in"}
> 设置 JUMP 模式下抬升高度和最大抬升高度 35](#_bookmark74)
>
> [![](media/image80.png){width="0.24166666666666667in" height="0.11in"}
> 获取 JUMP 模式下抬升高度和最大抬升高度 36](#_bookmark75)
>
> [![](media/image81.png){width="0.24166666666666667in"
> height="0.10999890638670166in"} 设置 JUMP 模式下扩展参数
> 36](#_bookmark76)
>
> [![](media/image82.png){width="0.238332239720035in"
> height="0.10999890638670166in"} 获取置 JUMP 模式下扩展参数
> 37](#_bookmark77)
>
> [![](media/image83.png){width="0.24166666666666667in"
> height="0.10999890638670166in"} 设置 PTP 模式下滑轨速度和加速度
> 37](#_bookmark78)
>
> [![](media/image84.png){width="0.315in"
> height="0.10999890638670166in"} 获取 PTP 模式下滑轨速度和加速度
> 38](#_bookmark79)
>
> [![](media/image85.png){width="0.29833333333333334in"
> height="0.10999890638670166in"} 设置 PTP
> 运动的速度百分比和加速度百分比 38](#_bookmark80)
>
> [![](media/image86.png){width="0.315in"
> height="0.10999890638670166in"} 获取 PTP
> 运动的速度百分比和加速度百分比 39](#_bookmark81)
>
> [![](media/image87.png){width="0.30833333333333335in"
> height="0.10999890638670166in"} 执行 PTP 指令 39](#_bookmark82)
>
> [![](media/image88.png){width="0.315in"
> height="0.10999890638670166in"} 执行带 I/O 控制的 PTP 指令
> 40](#_bookmark83)
>
> [![](media/image89.png){width="0.31166557305336834in"
> height="0.10999890638670166in"} 执行带滑轨的 PTP 指令
> 42](#_bookmark84)
>
> [![](media/image90.png){width="0.315in"
> height="0.10999890638670166in"} 执行带 I/O 控制和滑轨的 PTP 指令
> 43](#_bookmark85)
>
> [![](media/image91.png){width="0.198332239720035in"
> height="0.12333114610673666in"} **CP 功能** **46**](#_bookmark86)
>
> [![](media/image92.png){width="0.23in" height="0.10999890638670166in"}
> 设置 CP 运动的速度和加速度 46](#_bookmark87)
>
> [![](media/image93.png){width="0.24166666666666667in"
> height="0.10999890638670166in"} 获取 CP 运动的速度和加速度
> 46](#_bookmark88)
>
> [![](media/image94.png){width="0.235in"
> height="0.10999890638670166in"} 执行 CP 指令 47](#_bookmark89)
>
> [![](media/image95.png){width="0.24166666666666667in"
> height="0.10999890638670166in"} 执行带激光雕刻的 CP 指令
> 48](#_bookmark90)
>
> [![](media/image96.png){width="0.198332239720035in"
> height="0.12333114610673666in"} **ARC 功能 49**](#_bookmark91)
>
> [![](media/image97.png){width="0.23in" height="0.10999890638670166in"}
> 设置 ARC 运动的速度和加速度 49](#_bookmark92)
>
> [![](media/image98.png){width="0.24166666666666667in"
> height="0.10999890638670166in"} 获取 ARC 运动的速度和加速度
> 50](#_bookmark93)
>
> [![](media/image99.png){width="0.235in"
> height="0.10999890638670166in"} 执行 ARC 指令 50](#_bookmark94)
>
> [![](media/image100.png){width="0.24166666666666667in"
> height="0.10999890638670166in"} 执行 CIRCLE 指令 51](#_bookmark95)
>
> [![](media/image101.png){width="0.198332239720035in"
> height="0.12333114610673666in"} **丢步检测功能** **53**](#_bookmark96)
>
> [![](media/image102.png){width="0.23in"
> height="0.10999890638670166in"} 设置丢步检测阈值 53](#_bookmark97)
>
> [![](media/image103.png){width="0.24166666666666667in"
> height="0.10999890638670166in"} 执行丢步检测 53](#_bookmark98)
>
> [![](media/image104.png){width="0.235in"
> height="0.10999890638670166in"} 示例： 丢步检测 53](#_bookmark99)
>
> [![](media/image105.png){width="0.198332239720035in"
> height="0.12333333333333334in"} **WAIT 功能** **55**](#_bookmark100)
>
> [![](media/image106.png){width="0.23in"
> height="0.10999890638670166in"} 执行时间等待指令 55](#_bookmark101)
>
> [![](media/image107.png){width="0.24166666666666667in"
> height="0.10999890638670166in"} 执行触发指令 56](#_bookmark102)
>
> [![](media/image108.png){width="0.198332239720035in"
> height="0.12333114610673666in"} **EIO 功能** **57**](#_bookmark103)
>
> [![](media/image109.png){width="0.23in" height="0.11in"} 设置 I/O 复用
> 57](#_bookmark104)
>
> [![](media/image110.png){width="0.24166666666666667in"
> height="0.11in"} 读取 I/O 复用 58](#_bookmark105)
>
> [![](media/image111.png){width="0.235in" height="0.11in"} 设置 I/O
> 输出电平 58](#_bookmark106)
>
> [![](media/image112.png){width="0.24166666666666667in"
> height="0.11in"} 读取 I/O 输出电平 59](#_bookmark107)
>
> [![](media/image113.png){width="0.238332239720035in" height="0.11in"}
> 设置 PWM 输出 59](#_bookmark108)
>
> [![](media/image114.png){width="0.24166666666666667in"
> height="0.11in"} 读取 PWM 输出 60](#_bookmark109)
>
> [![](media/image115.png){width="0.24166666666666667in"
> height="0.11in"} 读取 I/O 输入电平 60](#_bookmark110)
>
> [![](media/image116.png){width="0.238332239720035in" height="0.11in"}
> 读取 A/D 输入 60](#_bookmark111)
>
> [![](media/image117.png){width="0.24166666666666667in"
> height="0.11in"} 设置扩展电机速度 61](#_bookmark112)
>
> [![](media/image118.png){width="0.315in" height="0.11in"}
> 设置扩展电机速度和移动距离 61](#_bookmark113)
>
> [![](media/image119.png){width="0.29833333333333334in"
> height="0.11in"} 使能光电传感器 62](#_bookmark114)
>
> [![](media/image120.png){width="0.315in" height="0.11in"}
> 获取光电传感器读数 62](#_bookmark115)
>
> [![](media/image121.png){width="0.30833333333333335in"
> height="0.11in"} 使能颜色传感器 63](#_bookmark116)
>
> [![](media/image122.png){width="0.315in" height="0.11in"}
> 获取颜色传感器读数 63](#_bookmark117)
>
> [![](media/image123.png){width="0.198332239720035in"
> height="0.12333114610673666in"} **CAL 功能** **65**](#_bookmark118)
>
> [![](media/image124.png){width="0.23in" height="0.11in"}
> 设置角度传感器静态偏差 65](#_bookmark119)
>
> [![](media/image125.png){width="0.24166666666666667in"
> height="0.10999890638670166in"} 读取角度传感器静态偏差
> 65](#_bookmark120)
>
> [![](media/image126.png){width="0.235in"
> height="0.10999890638670166in"} 设置角度传感器线性化参数
> 65](#_bookmark121)
>
> [![](media/image127.png){width="0.24166666666666667in"
> height="0.10999890638670166in"} 读取角度传感器线性化参数
> 66](#_bookmark122)
>
> [![](media/image128.png){width="0.238332239720035in"
> height="0.10999890638670166in"} 设置基座编码器静态偏差
> 66](#_bookmark123)
>
> [![](media/image129.png){width="0.24166666666666667in"
> height="0.10999890638670166in"} 读取基座编码器静态偏差
> 66](#_bookmark124)
>
> [![](media/image130.png){width="0.198332239720035in"
> height="0.12333114610673666in"} **WIFI 功能 67**](#_bookmark125)
>
> [![](media/image131.png){width="0.23in"
> height="0.10999890638670166in"} 使能 WIFI 67](#_bookmark126)
>
> [![](media/image132.png){width="0.24166666666666667in"
> height="0.10999890638670166in"} 获取 WIFI 状态 67](#_bookmark127)
>
> [![](media/image133.png){width="0.235in"
> height="0.10999890638670166in"} 设置 SSID 67](#_bookmark128)
>
> [![](media/image134.png){width="0.24166666666666667in"
> height="0.10999890638670166in"} 获取设置的 SSID 67](#_bookmark129)
>
> [![](media/image135.png){width="0.238332239720035in"
> height="0.10999890638670166in"} 设置 WIFI 密码 68](#_bookmark130)
>
> ![](media/image137.png){width="0.24499890638670166in"
> height="0.10999890638670166in"}[![](media/image138.png){width="0.24166666666666667in"
> height="0.10999890638670166in"} 获取 WIFI 密码 68](#_bookmark131)
>
> [![](media/image139.png){width="0.24166666666666667in"
> height="0.10999890638670166in"} 设置 IP 地址 68](#_bookmark132)
>
> [![](media/image140.png){width="0.238332239720035in"
> height="0.10999890638670166in"} 获取设置的 IP 地址 69](#_bookmark133)
>
> [![](media/image141.png){width="0.24166666666666667in"
> height="0.10999890638670166in"} 设置子网掩码 69](#_bookmark134)
>
> [![](media/image142.png){width="0.315in"
> height="0.10999890638670166in"} 获取设置的子网掩码 69](#_bookmark135)
>
> [![](media/image143.png){width="0.29833333333333334in"
> height="0.10999890638670166in"} 设置网关 70](#_bookmark136)
>
> [![](media/image144.png){width="0.315in"
> height="0.10999890638670166in"} 获取设置的网关 70](#_bookmark137)
>
> [![](media/image145.png){width="0.30833333333333335in"
> height="0.10999890638670166in"} 设置 DNS 70](#_bookmark138)
>
> [![](media/image146.png){width="0.315in"
> height="0.10999890638670166in"} 获取设置的 DNS 71](#_bookmark139)
>
> [![](media/image147.png){width="0.31166557305336834in"
> height="0.10999890638670166in"} 获取当前 WIFI 模块的连接状态
> 71](#_bookmark140)
>
> [![](media/image148.png){width="0.205in"
> height="0.12333333333333334in"} **其他功能** **72**](#_bookmark141)
>
> 事件循环功能 [72](#_bookmark118)
>
> []{#_bookmark1
> .anchor}![](media/image150.png){width="0.11833333333333333in"
> height="0.14333333333333334in"} **Dobot 指令简介**
>
> 控制器支持两种类型的指令： 立即指令与队列指令。
>
> . 立即指令：
> Dobot控制器在收到指令后立即处理该指令，而不管当前控制器是否在
> 还在处理其他指令。
>
> . 队列指令：
> Dobot控制器在收到指令后会将该指令放入控制器内部的指令队列中，
> Dobot控制器将顺序执行指令。
>
> 关于 Dobot 指令更具体的内容， 可查询 Dobot 通信协议手册。
>
> []{#_bookmark2
> .anchor}![](media/image152.png){width="0.13666557305336832in"
> height="0.14333333333333334in"} **指令超时**
>
> []{#_bookmark3 .anchor}![](media/image153.png){width="0.21in"
> height="0.12666557305336834in"} **设置指令超时时间**
>
> [如*1
> Dobot指令简介*中介绍，发送给Dobot控制器的所有指令都带有返回。当通信链路干](#_bookmark1)
> 扰等问题造成指令错误时，控制器将无法识别该条指令且无法返回。因此，每条下发给控制
> 器的指令都可设置一个超时时间。该指令超时时间可以通过以下的API进行设置。
>
> 表 2.1 设置指令超时时间

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 91%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetCmdTimeout(unsigned int cmdTimeout)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置指令超时时间。 当下发一条指令后如果需在规定时间内返回，则可调用
此接口设置超时时间，判断该指令是否超时</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>cmdTimeout：指令超时时间。单位： ms</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：无错误</p>
</blockquote></td>
</tr>
</tbody>
</table>

> 说明
>
> API返回结果为枚举类型，可在DobotType.h头文件中查看。
>
> []{#_bookmark4
> .anchor}![](media/image155.png){width="0.133332239720035in"
> height="0.14333333333333334in"} **连接/断开**
>
> []{#_bookmark5
> .anchor}![](media/image156.png){width="0.15833333333333333in"
> height="9.833333333333333e-2in"} **搜索** **Dobot**
>
> 表 3.1 搜索 Dobot

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 91%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SearchDobot(char *dobotNameList, uint32_tmaxLen)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>搜索Dobot，动态库将搜索到的Dobot信息存储，并使用ConnectDobot连接
Dobot</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>dobotNameList：字符串数组，Dobot动态库会将搜索到的串口或UDP信息写入
到dobotNameList 。一个典型的dobotNameList的格式为"COM1 COM3 COM6 <a
href="https://192.168.0.5">192.168.0.5</a>"
，多个串口或IP地址以空格分开</p>
<p>maxLen：字符串最大长度，以避免内存溢出</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>Dobot 数量</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark6
> .anchor}![](media/image157.png){width="0.22666557305336832in"
> height="0.12999890638670167in"} **连接** **Dobot**
>
> 表 3.2 连接 Dobot 控制器接口说明

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 91%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int ConnectDobot(const char *portName,uint32_t baudrate, char
*fwType, char *version, float *time)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>连接 Dobot。其中， portName 可以从 int SearchDobot(char *dobotList,
uint32_tmaxLen)的 char *dobotList 中得到</p>
<p>若 portName 为空， 并直接调用 ConnectDobot，则动态库将自动连接随机搜
索到的 Dobot</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>portName：Dobot 端口名。对于串口，其值为“COM3”。 对于 UDP，可能 是“
<a href="https://192.168.0.5">192.168.0.5</a> ”</p>
<p>baudrate：波特率，取值：115200</p>
<p>fwType：固件类型。包括：Dobot 和 Marlin</p>
<p>version：版本号</p>
<p>time：超时时间</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotConnect_NoError：连接 Dobot 成功</p>
<p>DobotConnect_NotFound：未找到 Dobot 端口</p>
<p>DobotConnect_Occupied：Dobot 端口被占用</p>
</blockquote></td>
</tr>
</tbody>
</table>

> ![](media/image158.png){width="0.2772594050743657in"
> height="0.23937992125984253in"}注意
>
> 请提前安装所需的驱动以便使API接口能识别Dobot控制器接口，详情请查询Dobot
> 用户手册。
>
> []{#_bookmark7
> .anchor}![](media/image160.png){width="0.22833333333333333in"
> height="0.12999890638670167in"} **断开** **Dobot**
>
> 表 3.3 断开 Dobot 接口说明

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 91%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int DisconnectDobot(void)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>断开 Dobot</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>无</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotConnect_NoError：无错误</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark8
> .anchor}![](media/image161.png){width="0.17499890638670165in"
> height="9.833333333333333e-2in"} **示例：** **连接示例**
>
> 程序 3.1 连接示例
>
> \#include \"DobotDll.h\"
>
> int split(char \*\*dst, char\* str, const char\* spl)
>
> {
>
> intn = 0;
>
> char \*result = NULL;
>
> result = strtok(str, spl);
>
> while( result != NULL )
>
> {
>
> strcpy(dst\[n++\], result);
>
> result = strtok(NULL, spl);
>
> }
>
> return n;
>
> }
>
> int main(void)
>
> {
>
> intmaxDevCount = 100;
>
> intmaxDevLen = 20;
>
> char \*devsChr = new char\[maxDevCount \* maxDevLen\]();
>
> char \*\*devsList = new char\*\[maxDevCount\]();
>
> for(int i=0; i\<maxDevCount; i++)
>
> devsList\[i\] = new char\[maxDevLen\]();
>
> SearchDobot(devsChr, 1024);
>
> split(devsList, devsChr, \" \");
>
> ConnectDobot(devsList\[0\], 115200, NULL, NULL, NULL);
>
> //控制 Dobot
>
> DisconnectDobot();
>
> delete\[\] devsChr;
>
> for(int i=0; i\<maxDevCount; i++)
>
> delete\[\] devsList\[i\];
>
> delete\[\] devsList;
>
> }
>
> []{#_bookmark9
> .anchor}![](media/image163.png){width="0.138332239720035in"
> height="0.14333333333333334in"} **指令队列控制**
>
> Dobot控制器中有一个存放指令的队列，以达到顺序存储和执行指令的目的。同时，通
> 过启动和停止指令， 还可实现丰富的异步操作。
>
> ![](media/image164.png){width="0.2772594050743657in"
> height="0.23937992125984253in"}注意
>
> 只有将"isQueued "参数设置为"1 "的指令才能加入指令队列。
>
> []{#_bookmark10
> .anchor}![](media/image165.png){width="0.21166557305336833in"
> height="0.12666557305336834in"} **执行队列中的指令**
>
> 表 4.1 执行指令

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 91%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetQueuedCmdStartExec(void)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>Dobot控制器开始循环查询指令队列，如果队列中有指令，则顺序取出并执行，
执行完一条指令后才会取出下一条继续执行</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>无</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark11
> .anchor}![](media/image166.png){width="0.22999890638670167in"
> height="0.12666557305336834in"} **停止执行队列中的指令**
>
> 表 4.2 停止执行指令

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 91%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetQueuedCmdStopExec(void)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>Dobot 控制器停止循环查询队列并停止执行指令。 在停止过程中若 Dobot 控
制器正在执行一条指令， 则待该指令执行完成后再停止。</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>无</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark12
> .anchor}![](media/image167.png){width="0.23166666666666666in"
> height="0.13in"} **强制停止执行队列中的指令**
>
> 表 4.3 强制停止执行指令

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 91%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetQueuedCmdForceStopExec(void)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>Dobot 控制器停止循环查询队列并停止执行指令。 在停止过程中若 Dobot 控
制器正在执行一条指令， 则该指令将会停止执行。</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>无</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark13
> .anchor}![](media/image169.png){width="0.23499890638670165in"
> height="0.12666557305336834in"} **示例：** **同步处理** **PTP
> 指令下发和队列控制**
>
> PTP详细[说明请参见12PTP功能。](#_bookmark69)
>
> 程序 4.1 同步处理 PTP 指令下发和队列控制
>
> \#include \"DobotDll.h\"
>
> int main(void)
>
> {
>
> uint64_t queuedCmdIndex = 0;
>
> PTPCmd cmd;
>
> cmd.ptpMode = 0;
>
> cmd.x
>
> cmd.y
>
> cmd.z
>
> cmd.r
>
> = 200;
>
> = 0;
>
> = 0;

= 0;

> ConnectDobot(NULL, 115200, NULL, NULL, NULL);
>
> SetQueuedCmdStartExec();
>
> SetPTPCmd(&cmd, true, &queuedCmdIndex);
>
> SetQueuedCmdStopExec();
>
> DisconnectDobot();
>
> }
>
> []{#_bookmark14
> .anchor}![](media/image170.png){width="0.23333333333333334in"
> height="0.13in"} **示例：** **异步处理** **PTP 指令下发和队列控制**
>
> 程序 4.2 异步处理 PTP 指令下发和队列控制
>
> int onButtonClick()
>
> {
>
> static bool flag = True;
>
> if (flag)
>
> SetQueuedCmdStartExec();
>
> else
>
> SetQueuedCmdStopExec();
>
> }
>
> // 子线程
>
> int thread(void)
>
> {
>
> uint64_t queuedCmdIndex = 0;
>
> PTPCmd cmd;
>
> cmd.ptpMode = 0;
>
> cmd.x
>
> cmd.y
>
> cmd.z
>
> cmd.r
>
> while(true)
>
> = 200;
>
> = 0;
>
> = 0;

= 0;

> SetPTPCmd(&cmd, true, &queuedCmdIndex);
>
> }
>
> []{#_bookmark15
> .anchor}![](media/image171.png){width="0.23333333333333334in"
> height="0.12999890638670167in"} **下载指令**
>
> Dobot 控制器支持将指令下载到控制器外部 Flash
> 中，然后通过控制器上的按键触发执 行，即脱机运行。
>
> 表 4.4 下载指令

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 91%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetQueuedCmdStartDownload(uint32_t totalLoop,uint32_t
linePerLoop)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>下载指令，当需要脱机运行时可调用此接口</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>totalLoop：脱机运行总次数</p>
<p>linePerLoop：指令单次循环次数。循环次数必须和下发的指令数相同，
且下发 指令必须设置为队列模式</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 91%" />
</colgroup>
<tbody>
<tr class="odd">
<td></td>
<td><blockquote>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark16
> .anchor}![](media/image172.png){width="0.23166666666666666in"
> height="0.12666447944007in"} **停止下载指令**
>
> 表 4.5 停止指令队列下载接口说明

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 91%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetQueuedCmdStopDownload(void)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>停止下载指令，当需要脱机运行时可调用此接口</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>无</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark17
> .anchor}![](media/image173.png){width="0.23166666666666666in"
> height="0.12999890638670167in"} **示例：** **下载** **PTP 指令**
>
> 程序 4.3 下载 PTP 指令
>
> \#include \"DobotDll.h\"
>
> int main(void)
>
> {
>
> uint64_t queuedCmdIndex = 0;
>
> PTPCmd cmd;
>
> cmd.ptpMode = 0;
>
> cmd.x
>
> cmd.y
>
> cmd.z
>
> cmd.r
>
> = 200;
>
> = 0;
>
> = 0;

= 0;

> ConnectDobot(NULL, 115200, NULL, NULL, NULL);
>
> //只下发一条 PTP 指令， linePerLoop 参数设为 1
>
> //总共循环两次，实际脱机运行两条 PTP 指令
>
> SetQueuedCmdStartDownload(2, 1);
>
> SetPTPCmd(&cmd, true, &queuedCmdIndex);
>
> SetQueuedCmdStopDownload();
>
> DisconnectDobot();
>
> }
>
> 指令下载的一般流程是：
>
> ![](media/image174.png){width="0.12666666666666668in"
> height="0.12333333333333334in"} 调用下载指令 API。
>
> ![](media/image175.png){width="0.14166557305336833in"
> height="0.12333333333333334in"} 发送指令并将指令设置为队列模式。
>
> ![](media/image176.png){width="0.138332239720035in"
> height="0.12333333333333334in"} 调用停止下载指令 API。
>
> []{#_bookmark18
> .anchor}![](media/image177.png){width="0.22999890638670167in"
> height="0.12999890638670167in"} **清空指令队列**
>
> 该接口可以清空Dobot控制器中的指令队列。
>
> 表 4.6 清除指令队列

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 91%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetQueuedCmdClear(void)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>清空指令队列</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>无</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark19
> .anchor}![](media/image178.png){width="0.323332239720035in"
> height="0.12999890638670167in"} **获取指令索引**
>
> 在 Dobot 控制器指令队列机制中，有一个 64
> 位内部计数器。当控制器每执行完一条指 令时，
> 该计数器将自动加一。通过该指令，可以查询当前执行完成的指令的索引。
>
> 表 4.7 获取指令队列当前索引接口说明

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 91%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetQueuedCmdCurrentIndex(uint64_t *queuedCmdCurrentIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取当前执行完成的指令的索引</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>queuedCmdCurrentIndex：指令索引</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark20
> .anchor}![](media/image179.png){width="0.3049989063867017in"
> height="0.12666557305336834in"} **示例：**
> **获取指令索引实现运动同步**
>
> 程序 4.4 获取指令索引实现运动同步
>
> \#include \"DobotDll.h\"
>
> int main(void)
>
> {
>
> uint64_t queuedCmdIndex = 0;
>
> uint64_t executedCmdIndex = 0;
>
> PTPCmd cmd;
>
> cmd.ptpMode = 0;
>
> cmd.x
>
> cmd.y
>
> cmd.z
>
> cmd.r
>
> = 200;
>
> = 0;
>
> = 0;

= 0;

> ConnectDobot(NULL, 115200, NULL, NULL, NULL);
>
> SetQueuedCmdStartExec();
>
> SetPTPCmd(&cmd, true, &queuedCmdIndex);
>
> // 通过比较队列索引实现运动同步
>
> While(executedCmdIndex \< queuedCmdIndex)
>
> GetQueuedCmdCurrentIndex(&executedCmdIndex);
>
> SetQueuedCmdStopExec();
>
> DisconnectDobot();
>
> }
>
> []{#_bookmark21
> .anchor}![](media/image181.png){width="0.13166666666666665in"
> height="0.14in"} **设备信息**
>
> []{#_bookmark22
> .anchor}![](media/image182.png){width="0.20666557305336833in"
> height="0.12999890638670167in"} **设置设备序列号**
>
> 表 5.1 设置设备序列号

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 91%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetDeviceSN(const char *deviceSN)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置设备序列号。该接口仅在出厂时有效（需要特殊密码）</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>deviceSN：字符串指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark23
> .anchor}![](media/image183.png){width="0.22499890638670167in"
> height="0.12999890638670167in"} **获取设备序列号**
>
> 表 5.2 获取设备序列号

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 91%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetDeviceSN(char *deviceSN, uint32_tmaxLen)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取设备序列号</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>deviceSN：字符串指针</p>
<p>maxLen：字符串最大长度，以避免溢出</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark24
> .anchor}![](media/image184.png){width="0.22666557305336832in"
> height="0.13in"} **设置设备名称**
>
> 表 5.3 设置设备名称

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 91%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetDeviceName(const char *deviceName)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置设备名称。当有多台机器时，可调用该接口设置设备名以作区分</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>deviceName：字符串指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark25
> .anchor}![](media/image185.png){width="0.22999890638670167in"
> height="0.13in"} **获取设备名称**
>
> 表 5.4 获取设备名称

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 91%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetDeviceName(char *deviceName,uint32_tmaxLen)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取设备名称。当有多台机器时，可调用该接口设置设备名以作区分</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 91%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>deviceName：字符串指针</p>
<p>maxLen：字符串最大长度，以避免溢出</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark26
> .anchor}![](media/image186.png){width="0.22833333333333333in"
> height="0.128332239720035in"} **获取设备版本号**
>
> 表 5.5 获取设备版本号

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 91%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetDeviceVersion(uint8_t *majorVersion,uint8_t
*minorVersion,uint8_t *revision)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取设备版本信息</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>majorVersion：主版本</p>
<p>minorVersion：次版本</p>
<p>revision：修订版本</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark27
> .anchor}![](media/image187.png){width="0.22833333333333333in"
> height="0.13in"} **设置滑轨状态**
>
> 表 5.6 设置滑轨状态

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 91%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetDeviceWithL(bool isEnable, bool isQueued,uint64_t
*queuedCmdIndex, uint8_t version)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置滑轨状态。若使用滑轨套件，则需调用该接口</p>
<p>在下发滑轨相关的指令前，必须先设置滑轨状态</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>isEnable：1 ，开启滑轨。 0 ，关闭滑轨</p>
<p>isQueued：是否将该指令加入指令队列</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。否
则，该参数无意义</p>
<p>version：滑轨版本号。0：版本为 V1.0 。1：版本为 V2.0</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark28
> .anchor}![](media/image188.png){width="0.22666557305336832in"
> height="0.128332239720035in"} **获取滑轨状态**
>
> 表 5.7 获取滑轨状态

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 91%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetDeviceWithL(bool *isEnable)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取滑轨状态</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>isEnable：1，开启。 0，关闭</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark29
> .anchor}![](media/image189.png){width="0.22666557305336832in"
> height="0.12999890638670167in"} **获取设备时钟**
>
> 表 5.8 获取设备时钟

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 91%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetDeviceTime(uint32_t *deviceTime)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取设备时钟</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>deviceTime：设备时钟</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark30
> .anchor}![](media/image191.png){width="0.13166666666666665in"
> height="0.14333333333333334in"} **实时位姿**
>
> 在 DobotV2.0 中， Dobot 控制器根据以下信息计算出实时位姿的基准值。
>
> . 底座码盘读数（可通过回零得到）。
>
> . 大臂角度传感器读数（上电或者按小臂UNLOCK按键时）。
>
> . 小臂角度传感器读数（上电或者按小臂UNLOCK按键时）。
>
> 在控制 Dobot 时， Dobot
> 控制器将基于实时位姿的基准值，以及实时运动状态，更新实 时位姿。
>
> []{#_bookmark31
> .anchor}![](media/image192.png){width="0.20666557305336833in"
> height="0.12999890638670167in"} **获取机械臂实时位姿**
>
> 表 6.1 获取实时位姿

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 11%" />
<col style="width: 17%" />
<col style="width: 61%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td colspan="3"><blockquote>
<p>int GetPose(Pose *pose)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td colspan="3"><blockquote>
<p>获取机械臂实时位姿</p>
</blockquote></td>
</tr>
<tr class="odd">
<td rowspan="3"><blockquote>
<p>参数</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>Pose 定义：</p>
<p>typedef struct tagPose {</p>
</blockquote></td>
<td rowspan="3"><blockquote>
<p>//机械臂坐标系 x</p>
<p>//机械臂坐标系 y</p>
<p>//机械臂坐标系 z</p>
<p>//机械臂坐标系 r</p>
<p>//机械臂关节轴(底座、大臂、小臂、末端)角度</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>float float float</p>
</blockquote>
<p>float</p></td>
<td><blockquote>
<p>x;</p>
<p>y;</p>
<p>z;</p>
<p>r;</p>
</blockquote></td>
</tr>
<tr class="odd">
<td colspan="2"><blockquote>
<p>float jointAngle[4];</p>
<p>}Pose;</p>
<p>Pose：Pose 指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td colspan="3"><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark32
> .anchor}![](media/image193.png){width="0.22499890638670167in"
> height="0.13in"} **获取滑轨实时位置**
>
> 表 6.2 获取滑轨实时位置

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetPoseL(flot *pose)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取导轨实时位姿</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>Pose：滑轨当前位置。单位 mm</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark33
> .anchor}![](media/image194.png){width="0.22666557305336832in"
> height="0.13in"} **重设机器人实时位姿**
>
> 在以下情况，需重新设置实时位姿的基准值：
>
> . 角度传感器损坏。
>
> . 角度传感器精度太差。
>
> 表 6.3 设置实时位姿的基准值

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int ResetPose(bool manual, float rearArmAngle, float
frontArmAngle)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>重新设置机械臂实时位姿</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>manual：表示是否自动重设姿态。0：自动重设姿态，无需设rearArmAngle及
frontArmAngle 。1：需设置rearArmAngle和frontArmAngle</p>
<p>rearArmAngle：大臂角度</p>
<p>frontArmAngle：小臂角度</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark34
> .anchor}![](media/image196.png){width="0.133332239720035in"
> height="0.14in"} **报警功能**
>
> []{#_bookmark35
> .anchor}![](media/image197.png){width="0.20666557305336833in"
> height="0.12666557305336834in"} **获取系统报警状态**
>
> 表 7.1 获取系统报警状态

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetAlarmsState(uint8_t *alarmsState,uint32_t *len, unsigned int
maxLen)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取系统报警状态</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>alarmsState：数组首地址。每一个字节可以标识8个报警项的报警状态，且MSB
（Most Significant Bit）在高位， LSB（Least Significant Bit）在低位</p>
<p>len：报警所占字节</p>
<p>maxLen：数组最大长度， 以避免溢出</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark36
> .anchor}![](media/image198.png){width="0.22499890638670167in"
> height="0.12666557305336834in"} **清除系统所有报警**
>
> 表 7.2 清除系统所有报警

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int ClearAllAlarmsState(void)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>清除系统所有报警</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>无</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark37
> .anchor}![](media/image200.png){width="0.128332239720035in"
> height="0.14333333333333334in"} **回零功能**
>
> 如果机械臂运行速度过快或者负载过大可能会导致精度降低，此时可执行回零操作，提
> 高精度。
>
> []{#_bookmark38
> .anchor}![](media/image201.png){width="0.20666557305336833in"
> height="0.12999890638670167in"} **设置回零位置**
>
> 表 8.1 设置回零位置

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetHOMEParams(HOMEParams *homeParams, bool isQueued,uint64_t
*queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置回零位置</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>HOMEParams 定义：</p>
<p>typedef struct tagHOMEParams {</p>
<p>float x; //机械臂坐标系 x</p>
<p>float y; //机械臂坐标系 y</p>
<p>float z; //机械臂坐标系 z</p>
<p>float r; //机械臂坐标系 r</p>
<p>}HOMEParams;</p>
<p>homeParams：HOMEParams 指针</p>
<p>isQueued：是否将该指令加入指令队列</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark39
> .anchor}![](media/image202.png){width="0.22499890638670167in"
> height="0.13in"} **获取回零位置**
>
> 表 8.2 获取回零位置

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 11%" />
<col style="width: 14%" />
<col style="width: 63%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td colspan="3"><blockquote>
<p>int GetHOMEParams(HOMEParams *homeParams)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td colspan="3"><blockquote>
<p>获取回零位置</p>
</blockquote></td>
</tr>
<tr class="odd">
<td rowspan="3"><blockquote>
<p>参数</p>
</blockquote></td>
<td colspan="3"><blockquote>
<p>HOMEParams 定义：</p>
<p>typedef struct tagHOMEParams {</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>float float float</p>
</blockquote>
<p>float</p></td>
<td><blockquote>
<p>x;</p>
<p>y;</p>
<p>z;</p>
<p>r;</p>
</blockquote></td>
<td><blockquote>
<p>//机械臂坐标系 x</p>
<p>//机械臂坐标系 y</p>
<p>//机械臂坐标系 z</p>
<p>//机械臂坐标系 r</p>
</blockquote></td>
</tr>
<tr class="odd">
<td colspan="3"><blockquote>
<p>}HOMEParams;</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td></td>
<td><blockquote>
<p>homeParams：HOMEParams 指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark40
> .anchor}![](media/image204.png){width="0.22666557305336832in"
> height="0.12999890638670167in"} **执行回零功能**
>
> 表 8.3 执行回零功能接口说明

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetHOMECmd(HOMECmd *homeCmd, bool isQueued,uint64_t
*queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>执行回零功能。调用该接口前如果未调用 SetHOMEParams 接口，则表示直
接回零至系统设置的位置。如果调用了 SetHOMEParams 接口，则回零至用
户自定义的位置</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>HOMECmd 定义：</p>
<p>typedef struct tagHOMECmd {</p>
<p>uint32_t reserved; //保留</p>
<p>}HOMECmd;</p>
<p>homeCmd：HOMECmd 指针</p>
<p>isQueued：是否将该指令加入指令队列</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark41
> .anchor}![](media/image205.png){width="0.22999890638670167in"
> height="0.13in"} **执行自动调平功能**
>
> 如果机械臂大小臂角度传感器出现偏差，导致机械臂定位精度降低，
> 可执行调平功能。 如果对定位精度要求较高，需手动调平，
> 详情请参见《Dobot Magicain 用户手册》。
>
> 表 8.4 执行自动调平功能接口说明

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetAutoLevelingCmd(AutoLevelingCmd *autoLevelingCmd, bool
isQueued, uint64_t *queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>执行自动调平功能</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>AutoLevelingCmd 定义：</p>
<p>typedef struct tagAutoLevelingCmd {</p>
</blockquote>
<table>
<colgroup>
<col style="width: 42%" />
<col style="width: 57%" />
</colgroup>
<tbody>
<tr class="odd">
<td>uint8_tcontrolFlag;</td>
<td><blockquote>
<p>//使能标志</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>float precision;</p>
</blockquote></td>
<td>//调平精度，最小值为 0.02</td>
</tr>
</tbody>
</table>
<blockquote>
<p>}AutoLevelingCmd;</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td></td>
<td><blockquote>
<p>autoLevelingCmd：AutoLevelingCmd 指针</p>
<p>isQueued：是否将该指令加入指令队列</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark42
> .anchor}![](media/image206.png){width="0.22833333333333333in"
> height="0.12999890638670167in"} **获取自动调平结果**
>
> 表 8.5 获取自动调平结果

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetAutoLevelingResult(float *precision)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取自动调平结果</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>Precision：精度</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark43
> .anchor}![](media/image208.png){width="0.133332239720035in"
> height="0.14333333333333334in"} **HHT 功能**
>
> HHT（Hand-hold Teaching）即手持示教。默认情况下，
> 用户按住小臂的圆形解锁按钮，
> 可拖动机械臂到任意位置，松开按钮就可以自动保存一个存点。
>
> []{#_bookmark44
> .anchor}![](media/image209.png){width="0.20833333333333334in"
> height="0.12999890638670167in"} **设置触发模式**
>
> 表 9.1 设置手持示教触发模式

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 42%" />
<col style="width: 48%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>int SetHHTTrigMode (HHTTrigMode hhtTrigMode)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>设置手持示教信号触发模式。如果不调用该函数， 则默认按键释放时触发</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>HHTTrigMode：</p>
<p>typedefenum tagHHTTrigMode{</p>
<p>TriggeredOnKeyReleased,</p>
<p>TriggeredOnPeriodicInterval</p>
<p>}HHTTrigMode;</p>
<p>hhtTrigMode：HHTTrigMode 枚举</p>
</blockquote></td>
<td><blockquote>
<p>//按键释放时触发</p>
<p>//按键被按下的过程中触发</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark45
> .anchor}![](media/image210.png){width="0.22666557305336832in"
> height="0.13in"} **获取触发模式**
>
> 表 9.2 获取手持示教触发模式

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 42%" />
<col style="width: 48%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>int GetHHTTrigMode (HHTTrigMode hhtTrigMode)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>获取手持示教信号触发模式</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>typedefenum tagHHTTrigMode{</p>
<p>TriggeredOnKeyReleased,</p>
<p>TriggeredOnPeriodicInterval</p>
<p>}HHTTrigMode;</p>
<p>hhtTrigMode：HHTTrigMode 枚举</p>
</blockquote></td>
<td><blockquote>
<p>//按键释放时触发</p>
<p>//按键被按下的过程中定时触发</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark46
> .anchor}![](media/image211.png){width="0.22833333333333333in"
> height="0.13in"} **设置手持示教使能状态**
>
> 表 9.3 设置手持示教使能状态

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetHHTTrigOutputEnabled (bool isEnabled)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置手持示教状态</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>isEnabled：0，去使能。 1，使能</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark47
> .anchor}![](media/image213.png){width="0.23166666666666666in"
> height="0.12999890638670167in"} **获取手持示教功能使能状态**
>
> 表 9.4 获取手持示教使能状态

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetHHTTrigOutputEnabled (bool *isEnabled)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取手持示教使能状态</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>isEnabled：0，去使能。 1，使能</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark48
> .anchor}![](media/image214.png){width="0.22999890638670167in"
> height="0.12999890638670167in"} **获取手持示教触发信号**
>
> 表 9.5 获取手持示教信号

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetHHTTrigOutput(bool *isTriggered)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取手持示教信号。需调用 SetHHTTrigOutputEnabled 接口后才能使用</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>isTriggered：0，未触发。1 ，已触发</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark49
> .anchor}![](media/image215.png){width="0.22999890638670167in"
> height="0.13in"} **示例：** **手持示教**
>
> 程序 9.1 手持示教
>
> \#include \"DobotDll.h\"
>
> \#include \<queue\>
>
> \#include \<windows.h\>
>
> int main(void)
>
> {
>
> ConnectDobot(NULL, 115200, NULL, NULL, NULL);
>
> SetHHTTrigMode(TriggeredOnPeriodicInterval);
>
> SetHHTTrigOutputEnabled(true);
>
> bool isTriggered = false;
>
> []{#_bookmark50
> .anchor}![](media/image217.png){width="0.21666557305336834in"
> height="0.14333333333333334in"} **末端执行器**
>
> []{#_bookmark51
> .anchor}![](media/image218.png){width="0.2949989063867017in"
> height="0.12999890638670167in"} **设置末端坐标偏移量**
>
> 表 10.1 设置末端坐标偏移量

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 11%" />
<col style="width: 32%" />
<col style="width: 45%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td colspan="3"><blockquote>
<p>int SetEndEffectorParams(EndEffectorParams *endEffectorParams, bool
isQueued, uint64_t *queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td colspan="3"><blockquote>
<p>设置末端坐标偏移参数， 一般在末端安装了执行器才需设置</p>
<p>当使用标准的末端执行器（机械臂配套的末端执行器）时，请查询 Dobot 用
户手册，得到其 X 轴与 Y 轴偏置，并调用本接口设置。其他情况下的末端
执行器参数，需自行确认结构参数。</p>
</blockquote></td>
</tr>
<tr class="odd">
<td rowspan="4"><blockquote>
<p>参数</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>EndEffectorParams 定义：</p>
<p>typedef struct tagEndEffectorParams {</p>
</blockquote></td>
<td rowspan="3"><blockquote>
<p>//末端 X 方向偏移量</p>
<p>//末端 Y 方向偏移量</p>
<p>//末端 Z 方向偏移量</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>float float</p>
</blockquote>
<p>float</p></td>
<td><blockquote>
<p>xBias;</p>
<p>yBias;</p>
<p>zBias;</p>
</blockquote></td>
</tr>
<tr class="odd">
<td colspan="2"><blockquote>
<p>}EndEffectorParams;</p>
</blockquote></td>
</tr>
<tr class="even">
<td colspan="3"><blockquote>
<p>endEffectorParams：EndEffectorParams 指针</p>
<p>isQueued：是否将该指令加入指令队列</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td colspan="3"><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark52
> .anchor}![](media/image219.png){width="0.313332239720035in"
> height="0.13in"} **获取末端坐标偏移量**
>
> 表 10.2 获取末端坐标偏移量

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 11%" />
<col style="width: 32%" />
<col style="width: 45%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td colspan="3"><blockquote>
<p>int GetEndEffectorParams(EndEffectorParams *endEffectorParams)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td colspan="3"><blockquote>
<p>获取末端坐标偏移量</p>
</blockquote></td>
</tr>
<tr class="odd">
<td rowspan="4"><blockquote>
<p>参数</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>EndEffectorParams 定义：</p>
<p>typedef struct tagEndEffectorParams {</p>
</blockquote></td>
<td rowspan="3"><blockquote>
<p>//末端 X 方向偏移量</p>
<p>//末端 Y 方向偏移量</p>
<p>//末端 Z 方向偏移量</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>float float</p>
</blockquote>
<p>float</p></td>
<td><blockquote>
<p>xBias;</p>
<p>yBias;</p>
<p>zBias;</p>
</blockquote></td>
</tr>
<tr class="odd">
<td colspan="2"><blockquote>
<p>}EndEffectorParams;</p>
</blockquote></td>
</tr>
<tr class="even">
<td colspan="3"><blockquote>
<p>endEffectorParams：EndEffectorParams 指针</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark53 .anchor}![](media/image220.png){width="0.315in"
> height="0.12999890638670167in"} **设置激光状态**
>
> 表 10.3 设置激光状态

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetEndEffectorLaser(bool enableCtrl, bool on, bool
isQueued,uint64_t *queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置激光状态</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>enableCtrl：末端使能。 0，未使能。1，使能</p>
<p>on：开启或停止激光。0，停止。1，开启</p>
<p>isQueued：是否将该指令加入指令队列</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark54
> .anchor}![](media/image221.png){width="0.318332239720035in"
> height="0.13in"} **获取激光状态**
>
> 表 10.4 获取激光状态

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetEndEffectorLaser(bool *isCtrlEnabled, bool *isOn)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取激光状态</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>isCtrlEnabled：末端是否使能。0，未使能. 1，使能</p>
<p>isOn：激光是否开启。0，停止. 1，开启</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark55
> .anchor}![](media/image222.png){width="0.31666557305336834in"
> height="0.13in"} **设置气泵状态**
>
> 表 10.5 设置气泵状态

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetEndEffectorSuctionCup(bool enableCtrl, bool suck, bool
isQueued, uint64_t *queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置气泵状态</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>enableCtrl：末端使能。 0，未使能. 1，使能</p>
<p>suck：控制气泵吸气或吹气。 0：吹气。1，吸气</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td></td>
<td><blockquote>
<p>isQueued：是否将该指令加入指令队列</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark56
> .anchor}![](media/image223.png){width="0.31666557305336834in"
> height="0.12999890638670167in"} **获取气泵状态**
>
> 表 10.6 获取气泵状态

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetEndEffectorSuctionCup(bool *isCtrlEnabled, bool *isSucked)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取气泵状态</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>isCtrlEnabled：末端是否使能。 0，未使能。1，使能</p>
<p>isSucked:气泵是否吸气。 0：吹气。1，吸气</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark57 .anchor}![](media/image224.png){width="0.315in"
> height="0.13in"} **设置夹爪状态**
>
> 表 10.7 设置夹爪状态

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetEndEffectorGripper(bool enableCtrl, bool grip, bool
isQueued,uint64_t *queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置夹爪状态</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>enableCtrl：末端使能。 0，未使能。1，使能</p>
<p>grip：控制夹爪抓取或释放。 0，释放。1，抓取</p>
<p>isQueued：是否将该指令加入指令队列</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark58 .anchor}![](media/image225.png){width="0.315in"
> height="0.13in"} **获取夹爪状态**
>
> 表 10.8 获取状态接口说明

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetEndEffectorGripper(bool *isCtrlEnabled, bool *isGripped)</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取夹爪状态</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>isCtrlEnabled：末端是否使能。 0，未使能。1，使能</p>
<p>isGripped：夹爪是否抓取。0，释放。1，抓取</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark59
> .anchor}![](media/image227.png){width="0.21666557305336834in"
> height="0.14333333333333334in"} **JOG 功能**
>
> 机械臂点动时各轴实际运动速度与设置的速度关系如下所示：
>
> . 各轴点动速度=各轴设置的点动速度\*设置的速度百分比
>
> . 各轴点动加速度=各轴设置的点动加速度\*设置的加速度百分比
>
> []{#_bookmark60
> .anchor}![](media/image228.png){width="0.2949989063867017in"
> height="0.12666557305336834in"}
> **设置点动时各关节坐标轴的动速度和加速度**
>
> 表 11.1 设置点动时各关节坐标轴的速度和加速度

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetJOGJointParams(JOGJointParams *jogJointParams, bool isQueued,
uint64_t *queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置点动时各关节坐标轴的速度( °/s）和加速度( °/s2）</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>JOGJointParams 定义：</p>
<p>typedef struct tagJOGJointParams {</p>
</blockquote>
<table>
<colgroup>
<col style="width: 52%" />
<col style="width: 47%" />
</colgroup>
<tbody>
<tr class="odd">
<td>float velocity[4];</td>
<td><blockquote>
<p>//4 轴关节速度</p>
</blockquote></td>
</tr>
<tr class="even">
<td>float acceleration[4];</td>
<td>//4 轴关节加速度</td>
</tr>
</tbody>
</table>
<blockquote>
<p>}JOGJointParams;</p>
<p>jogJointParam：JOGJointParams 指针</p>
<p>isQueued：是否将该指令加入指令队列</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark61
> .anchor}![](media/image229.png){width="0.313332239720035in"
> height="0.12666557305336834in"}
> **获取点动时各关节坐标轴的动速度和加速度**
>
> 表 11.2 获取点动时各关节坐标轴的速度和加速度

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetJOGJointParams(JOGJointParams *jogJointParams)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取点动时各关节坐标轴的速度( °/s）和加速度( °/s2）</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>JOGJointParams 定义：</p>
<p>typedef struct tagJOGJointParams {</p>
</blockquote>
<table>
<colgroup>
<col style="width: 52%" />
<col style="width: 47%" />
</colgroup>
<tbody>
<tr class="odd">
<td>float velocity[4];</td>
<td><blockquote>
<p>//4 轴关节速度</p>
</blockquote></td>
</tr>
<tr class="even">
<td>float acceleration[4];</td>
<td>//4 轴关节加速度</td>
</tr>
</tbody>
</table>
<blockquote>
<p>}JOGJointParams</p>
<p>jogJointParams：JOGJointParams 指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td></td>
<td><blockquote>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark62 .anchor}![](media/image230.png){width="0.315in"
> height="0.12999890638670167in"}
> **设置点动时笛卡尔坐标轴的速度和加速度**
>
> 表 11.3 设置点动时笛卡尔坐标轴的速度和加速度

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetJOGCoordinateParams(JOGCoordinateParams
*jogCoordinateParams,</p>
<p>bool isQueued,uint64_t *queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置点动时各坐标轴（笛卡尔） 的速度（mm/s）和加速度（mm/s2）</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>JOGCoordinateParams 定义：</p>
<p>typedef struct tagJOGCoordinateParams {</p>
<p>float velocity[4]; //4 轴坐标轴 X,Y,Z,R 速度</p>
<p>float acceleration[4]; //4 轴坐标轴 X,Y,Z,R 加速度</p>
<p>}JOGCoordinateParams;</p>
<p>jogCoordinateParams：JOGCoordinateParams 指针</p>
<p>isQueued：是否将该指令加入指令队列</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark63
> .anchor}![](media/image231.png){width="0.318332239720035in"
> height="0.12666557305336834in"}
> **获取点动时笛卡尔坐标轴的速度和加速度**
>
> 表 11.4 获取点动时笛卡尔坐标轴的速度和加速度

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetJOGCoordinateParams(JOGCoordinateParams
*jogCoordinateParams)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取点动时各坐标轴（笛卡尔） 的速度（mm/s）和加速度（mm/s2）</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>JOGCoordinateParams 定义：</p>
<p>typedef struct tagJOGCoordinateParams {</p>
<p>float velocity[4]; //4 轴坐标轴 X,Y,Z,R 速度</p>
<p>float acceleration[4]; //4 轴坐标轴 X,Y,Z,R 加速度</p>
<p>}JOGCoordinateParams;</p>
<p>jogCoordinateParams： JOGCoordinateParams 指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark64
> .anchor}![](media/image232.png){width="0.31666557305336834in"
> height="0.13in"} **设置点动时滑轨速度和加速度**
>
> 表 11.5 设置点动时滑轨速度和加速度

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetJOGLParams(JOGLParams *jogLParams, bool isQueued,uint64_t
*queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置点动时滑轨的速度（mm/s）和加速度（mm/s2）</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>jogLParams 定义：</p>
<p>typedef struct tagJogLParams {</p>
<p>float velocity; //滑轨速度</p>
<p>float acceleration; //滑轨加速度</p>
<p>}JogLParams;</p>
<p>jogLParams：jogLParams 指针</p>
<p>isQueued：是否将该指令加入指令队列</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark65
> .anchor}![](media/image233.png){width="0.31666557305336834in"
> height="0.13in"} **获取点动时滑轨速度和加速度**
>
> 表 11.6 获取点动时滑轨速度和加速度

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetJOGLParams(JOGLParams * jogLParams )</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取点动时滑轨的速度（mm/s）和加速度（mm/s2）</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>jogLParams 定义：</p>
<p>typedef struct tagJogLParams {</p>
<p>float velocity; //滑轨速度</p>
<p>float acceleration; //滑轨加速度</p>
<p>}JogLParams;</p>
<p>jogLParams：jogLParams 指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark66 .anchor}![](media/image234.png){width="0.315in"
> height="0.12666557305336834in"} **设置点动速度百分比和加速度百分比**
>
> 表 11.7 设置点动速度百分比和加速度百分比

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetJOGCommonParams(JOGCommonParams *jogCommonParams, bool
isQueued, uint64_t *queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置点动（关节坐标系下和笛卡尔坐标系下） 速度百分比和加速度百分比</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>JOGCommonParams 定义：</p>
<p>typedef struct tagJOGCommonParams {</p>
<p>float velocityRatio; //速度比例，关节坐标轴点动和笛卡尔坐标轴</p>
<p>点动共用</p>
<p>float accelerationRatio; //加速度比例，关节坐标轴点动和笛卡尔坐标轴
点动共用</p>
<p>}JOGCommonParams;</p>
<p>jogCommonParams：JOGCommonParams 指针</p>
<p>isQueued：是否将该指令加入指令队列</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark67 .anchor}![](media/image235.png){width="0.315in"
> height="0.12999890638670167in"} **获取点动速度百分比和加速度百分比**
>
> 表 11.8 获取点动速度百分比和加速度百分比

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetJOGCommonParams(JOGCommonParams *jogCommonParams)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取点动（关节坐标系下和笛卡尔坐标系下） 速度百分比和加速度百分比</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>JOGCommonParams 定义：</p>
<p>typedef struct tagJOGCommonParams {</p>
<p>float velocityRatio; //速度比例，关节点动和坐标轴点动共用</p>
<p>float accelerationRatio; //加速度比例，关节点动和坐标轴点动共用</p>
<p>}JOGCommonParams;</p>
<p>jogCommonParams：JOGCommonParams 指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark68
> .anchor}![](media/image236.png){width="0.313332239720035in"
> height="0.13in"} **执行点动指令**
>
> 表 11.9 执行点动指令接口说明

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetJOGCmd(JOGCmd *jogCmd, bool isQueued,uint64_t
*queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>执行点动指令。设置点动相关参数后可调用该接口</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>JOGCmd 定义：</p>
<p>typedef struct tagJOGCmd {</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 26%" />
<col style="width: 64%" />
</colgroup>
<tbody>
<tr class="odd">
<td rowspan="2"></td>
<td><blockquote>
<p>uint8_t isJoint;</p>
<p>uint8_t cmd;</p>
<p>}JOGCmd;</p>
<p>//点动命令详细说明</p>
<p>enum {</p>
<p>IDLE,</p>
<p>AP_DOWN,</p>
<p>AN_DOWN,</p>
<p>BP_DOWN,</p>
<p>BN_DOWN,</p>
<p>CP_DOWN,</p>
<p>CN_DOWN,</p>
<p>DP_DOWN,</p>
<p>DN_DOWN,</p>
<p>LP_DOWN,</p>
<p>LN_DOWN</p>
<p>};</p>
</blockquote></td>
<td><blockquote>
<p>//点动方式：0 ，笛卡尔坐标轴点动；1 ，关节点动 //点动命令（取值范围
0~10）</p>
<p>//空闲状态</p>
<p>//X+/Joint1+</p>
<p>//X-/Joint1-</p>
<p>//Y+/Joint2+</p>
<p>//Y-/Joint2-</p>
<p>//Z+/Joint3+</p>
<p>//Z-/Joint3-</p>
<p>//R+/Joint4+</p>
<p>//R-/Joint4-</p>
<p>//L+ 。仅在 isJoint=1 时， LP_DOWN 可用</p>
<p>//L- 。仅在 isJoint=1 时， LN_DOWN 可用</p>
</blockquote></td>
</tr>
<tr class="even">
<td colspan="2"><blockquote>
<p>jogCmd：JOGCmd 指针</p>
<p>isQueued：是否将该指令加入指令队列</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark69
> .anchor}![](media/image238.png){width="0.21666557305336834in"
> height="0.14333333333333334in"} **PTP 功能**
>
> PTP点位模式点位模式即实现点到点运动，Dobot M1的点位模式包括MOVJ
> 、MOVL以 及JUMP三种运动模式。不同的运动模式，
> 示教后存点回放的运动轨迹不同。
>
> . MOVJ：关节运动，由A点运动到B点，各个关节从A点对应的关节角运行至B点对
> 应的关节角。关节运动过程中，
> 各个关节轴的运行时间需一致，且同时到达终点， [如图
> 12.1所](#_bookmark144)示。

![](media/image239.jpeg){width="2.875in" height="1.9166666666666667in"}

> []{#_bookmark144 .anchor}图 12.1 MOVL 和 MOVJ 运动模式

.

.

MOVL：直线运动， A点到B[点的路径为直线，如图 12.1所](#_bookmark144)示。

JUMP：门型轨迹[，如图 12.2所](#_bookmark70)示， 由A点到B点的JUMP运动，
先抬升高度Height， 再平移到B点上方Height处，然后下降Height。

![](media/image240.jpeg){width="3.7766655730533683in"
height="2.028332239720035in"}

> 图 12.2 JUMP 运动模式
>
> 机械臂再现运动时各轴实际运动速度与设置的速度关系如下所示：
>
> . 各轴再现速度=各轴设置的再现速度\*设置的速度百分比
>
> . 各轴再现加速度=各轴设置的再现加速度\*设置的加速度百分比
>
> []{#_bookmark70
> .anchor}![](media/image241.png){width="0.2949989063867017in"
> height="0.12666557305336834in"} **设置** **PTP
> 模式下各关节坐标轴的速度和加速度**
>
> 表 12.1 设置 PTP 模式下各关节坐标轴的速度和加速度

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetPTPJointParams(PTPJointParams *ptpJointParams, bool
isQueued,</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td></td>
<td><blockquote>
<p>uint64_t *queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置 PTP 运动时各关节坐标轴的速度( °/s）和加速度( °/s2）</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>PTPJointParams 定义：</p>
<p>typedef struct tagPTPJointParams{</p>
<p>float velocity[4]; //PTP 模式下 4 轴关节速度</p>
<p>float acceleration[4]; //PTP 模式下 4 轴关节加速度</p>
<p>}PTPJointParams;</p>
<p>ptpJointParams：PTPJointParams 指针</p>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark71
> .anchor}![](media/image243.png){width="0.313332239720035in"
> height="0.12666557305336834in"} **获取** **PTP
> 模式下各关节坐标轴的速度和加速度**
>
> 表 12.2 获取关节点位参数接口说明

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetPTPJointParams(PTPJointParams *ptpJointParams)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取 PTP 运动时各关节坐标轴的速度( °/s）和加速度( °/s2）</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>PTPJointParams 定义：</p>
<p>typedef struct tagPTPJointParams {</p>
<p>float velocity[4]; //PTP 模式下 4 轴关节速度</p>
<p>float acceleration[4]; //PTP 模式下 4 轴关节加速度</p>
<p>}PTPJointParams;</p>
<p>ptpJointParams：PTPJointParams 指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark72 .anchor}![](media/image244.png){width="0.315in"
> height="0.13in"} **设置** **PTP 模式下各笛卡尔坐标轴的速度和加速度**
>
> 表 12.3 设置 PTP 运动时各笛卡尔关节坐标轴的速度和加速度

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetPTPCoordinateParams(PTPCoordinateParams *ptpCoordinateParams,
bool isQueued, uint64_t *queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置 PTP 运动时各笛卡尔坐标轴的速度（mm/s）和加速度（mm/s2）</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>PTPCoordinateParams 定义：</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td></td>
<td><blockquote>
<p>typedef struct tagPTPCoordinateParams {</p>
<p>float xyzVelocity; //PTP 模式下 X,Y,Z 3 轴坐标轴速度</p>
<p>float rVelocity; //PTP 模式下末端 R 轴速度</p>
<p>float xyzAcceleration; //PTP 模式下 X,Y,Z 3 轴坐标轴加速度</p>
<p>float rAccleration; //PTP 模式下末端 R 轴加速度</p>
<p>} PTPCoordinateParams;</p>
<p>ptpCoordinateParams：PTPCoordinateParams 指针</p>
<p>isQueued:是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark73
> .anchor}![](media/image245.png){width="0.318332239720035in"
> height="0.12666557305336834in"} **获取** **PTP
> 模式下各笛卡尔坐标轴的速度和加速度**
>
> 表 12.4 获取 PTP 运动时各笛卡尔关节坐标轴的速度 s 和加速度

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetPTPCoordinateParams(PTPCoordinateParams
*ptpCoordinateParams)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>PTP 运动时各笛卡尔坐标轴的速度（mm/s）和加速度（mm/s2）</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>PTPCoordinateParams 定义：</p>
<p>typedef struct tagPTPCoordinateParams {</p>
<p>float xyzVelocity; //PTP 模式下 X,Y,Z 3 轴坐标轴速度</p>
<p>float rVelocity; //PTP 模式下末端 R 轴速度</p>
<p>float xyzAcceleration; //PTP 模式下 X,Y,Z 3 轴坐标轴加速度</p>
<p>float rAccleration; //PTP 模式下末端 R 轴加速度</p>
<p>} PTPCoordinateParams;</p>
<p>ptpCoordinateParams：PTPCoordinateParams 指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark74
> .anchor}![](media/image246.png){width="0.31666557305336834in"
> height="0.13in"} **设置** **JUMP 模式下抬升高度和最大抬升高度**
>
> 表 12.5 设置 JUMP 模式下抬升高度和最大抬升高度

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetPTPJumpParams(PTPJumpParams *ptpJumpParams, bool isQueued,
uint64_t *queuedCmdIndex)</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置 JUMP 模式下抬升高度和最大抬升高度</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>PTPJumpParams 定义：</p>
<p>typedef struct tagPTPJumpParams {</p>
<p>float jumpHeight; //抬升高度</p>
<p>float zLimit; //最大抬升高度</p>
<p>}PTPJumpParams;</p>
<p>ptpJumpParams：PTPJumpParams 指针</p>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark75
> .anchor}![](media/image247.png){width="0.31666557305336834in"
> height="0.13in"} **获取** **JUMP 模式下抬升高度和最大抬升高度**
>
> 表 12.6 JUMP 模式下抬升高度和最大抬升高度

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetPTPJumpParams(PTPJumpParams *ptpJumpParams)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取 JUMP 模式下抬升高度和最大抬升高度</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>PTPJumpParams 定义：</p>
<p>typedef struct tagPTPJumpParams {</p>
<p>float jumpHeight; //抬升高度</p>
<p>float zLimit; //最大抬升高度</p>
<p>}PTPJumpParams;</p>
<p>ptpJumpParams：PTPJumpParams 指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark76 .anchor}![](media/image248.png){width="0.315in"
> height="0.12666557305336834in"} **设置** **JUMP 模式下扩展参数**
>
> 表 12.7 设置 JUMP 模式下扩展参数

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetPTPJump2Params(PTPJumpParams *ptpJump2Params, bool isQueued,
uint64_t *queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置 JUMP 模式下扩展参数</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>PTPJump2Params 定义：</p>
<p>typedef struct tagPTPJump2Params {</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td></td>
<td><table>
<colgroup>
<col style="width: 52%" />
<col style="width: 47%" />
</colgroup>
<tbody>
<tr class="odd">
<td>float startJumpHeight;</td>
<td><blockquote>
<p>//起始点抬升高度</p>
</blockquote></td>
</tr>
<tr class="even">
<td>float endJumpHeight;</td>
<td>//结束点抬升高度</td>
</tr>
<tr class="odd">
<td>float zLimit;</td>
<td><blockquote>
<p>//最大抬升高度</p>
</blockquote></td>
</tr>
</tbody>
</table>
<blockquote>
<p>}PTPJump2Params;</p>
<p>ptpJump2Params：PTPJump2Params 指针</p>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark77 .anchor}![](media/image249.png){width="0.315in"
> height="0.12999890638670167in"} **获取置** **JUMP 模式下扩展参数**
>
> 表 12.8 获取 JUMP 模式下扩展参数

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetPTPJump2Params(PTPJumpParams *ptpJump2Params)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取 JUMP 模式下扩展参数</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>PTPJump2Params 定义：</p>
<p>typedef struct tagPTPJump2Params {</p>
</blockquote>
<table>
<colgroup>
<col style="width: 52%" />
<col style="width: 47%" />
</colgroup>
<tbody>
<tr class="odd">
<td>float startJumpHeight;</td>
<td><blockquote>
<p>//起始点抬升高度</p>
</blockquote></td>
</tr>
<tr class="even">
<td>float endJumpHeight;</td>
<td>//结束点抬升高度</td>
</tr>
<tr class="odd">
<td>float zLimit;</td>
<td><blockquote>
<p>//最大抬升高度</p>
</blockquote></td>
</tr>
</tbody>
</table>
<blockquote>
<p>}PTPJump2Params;</p>
<p>ptpJump2Params：PTPJump2Params 指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark78
> .anchor}![](media/image250.png){width="0.313332239720035in"
> height="0.13in"} **设置** **PTP 模式下滑轨速度和加速度**
>
> 表 12.9 设置 PTP 模式下滑轨速度和加速度

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetPTPLParams(PTPLParams * ptpLParams,bool isQueued,uint64_t
*queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置 PTP 模式下滑轨速度（mm/s）和加速度（mm/s2）</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>PTPLParams 定义：</p>
<p>typedef struct tagPTPLParams {</p>
<p>float velocity; //滑轨速度</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td></td>
<td><blockquote>
<p>float acceleration; //滑轨加速度</p>
<p>}PTPLParams;</p>
<p>ptpLParams：PTPLParams 指针</p>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark79
> .anchor}![](media/image251.png){width="0.40499890638670166in"
> height="0.12999890638670167in"} **获取** **PTP
> 模式下滑轨速度和加速度**
>
> 表 12.10 获取 PTP 模式下滑轨速度和加速度

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetPTPLParams(PTPLParams *ptpLParams)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取 PTP 模式下滑轨速度（mm/s）和加速度（mm/s2）</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>PTPLParams 定义：</p>
<p>typedef struct tagPTPLParams {</p>
<p>float velocity; //滑轨速度</p>
<p>float acceleration; //滑轨加速度</p>
<p>}PTPLParams;</p>
<p>ptpLParams：PTPLParams 指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark80
> .anchor}![](media/image252.png){width="0.38666666666666666in"
> height="0.12666557305336834in"} **设置** **PTP
> 运动的速度百分比和加速度百分比**
>
> 表 12.11 设置 PTP 运动的速度百分比和加速度百分比

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 11%" />
<col style="width: 20%" />
<col style="width: 58%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td colspan="3"><blockquote>
<p>int SetPTPCommonParams(PTPCommonParams *ptpCommonParams, bool
isQueued, uint64_t *queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td colspan="3"><blockquote>
<p>设置 PTP 运动的速度百分比和加速度百分比</p>
</blockquote></td>
</tr>
<tr class="odd">
<td rowspan="3"><blockquote>
<p>参数</p>
</blockquote></td>
<td colspan="3"><blockquote>
<p>PTPCommonParams 定义：</p>
<p>typedef struct tagPTPCommonParams {</p>
</blockquote></td>
</tr>
<tr class="even">
<td><p>float</p>
<p>float</p></td>
<td><blockquote>
<p>velocityRatio;</p>
<p>accelerationRatio;</p>
</blockquote></td>
<td><blockquote>
<p>//PTP 模式速度百分比，关节坐标轴和笛卡尔 坐标轴共用</p>
<p>//PTP 模式加速度百分比，关节坐标轴和笛卡尔 坐标轴共用</p>
</blockquote></td>
</tr>
<tr class="odd">
<td colspan="3"><blockquote>
<p>}PTPCommonParams;</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 8%" />
<col style="width: 91%" />
</colgroup>
<tbody>
<tr class="odd">
<td></td>
<td><blockquote>
<p>ptpCommonParams：PTPCommonParams 指针</p>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。否
则，该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark81
> .anchor}![](media/image253.png){width="0.40499890638670166in"
> height="0.128332239720035in"} **获取** **PTP
> 运动的速度百分比和加速度百分比**
>
> 表 12.12 PTP 运动的速度百分比和加速度百分比

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetPTPCommonParams(PTPCommonParams *ptpCommonParams)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>PTP 运动的速度百分比和加速度百分比</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>PTPCommonParams 定义：</p>
<p>typedef struct tagPTPCommonParams {</p>
<p>float velocityRatio; //PTP 模式速度百分比，关节坐标轴和笛卡尔</p>
<p>坐标轴共用</p>
<p>float accelerationRatio; //PTP 模式加速度百分比，关节坐标轴和笛卡尔
坐标轴共用</p>
<p>}PTPCommonParams;</p>
<p>ptpCommonParams：PTPCommonParams 指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark82
> .anchor}![](media/image254.png){width="0.4066655730533683in"
> height="0.13in"} **执行** **PTP 指令**
>
> 表 12.13 执行 PTP 指令

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetPTPCmd(PTPCmd *ptpCmd, bool isQueued,uint64_t
*queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>执行 PTP 指令。设置 PTP 相关参数后，调用此函数可使机械臂运动至设置的
目标点。</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>PTPCmd 定义：</p>
<p>typedef struct tagPTPCmd {</p>
<p>uint8_tptpMode; //PTP 模式，取值范围：0~9</p>
<p>float x; //（x,y,z,r）为坐标参数，可为笛卡尔坐标、关节坐</p>
<p>标、 笛卡尔坐标增量或关节坐标增量</p>
<p>float y;</p>
<p>float z;</p>
<p>float r;</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 27%" />
<col style="width: 63%" />
</colgroup>
<tbody>
<tr class="odd">
<td rowspan="3"></td>
<td colspan="2"><blockquote>
<p>}PTPCmd;</p>
<p>//其中， ptpMode 取值如下：</p>
<p>enum {</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>JUMP_XYZ,</p>
<p>MOVJ_XYZ,</p>
<p>MOVL_XYZ,</p>
<p>JUMP_ANGLE,</p>
<p>MOVJ_ANGLE,</p>
<p>MOVL_ANGLE,</p>
<p>MOVJ_INC,</p>
<p>MOVL_INC,</p>
</blockquote></td>
<td><blockquote>
<p>//JUMP 模式，（x,y,z,r）为笛卡尔坐标系下的目标 点坐标</p>
<p>//MOVJ 模式，（x,y,z,r）为笛卡尔坐标系下的目标 点坐标</p>
<p>//MOVL 模式，（x,y,z,r）为笛卡尔坐标系下的目标 点坐标</p>
<p>//JUMP 模式，（x,y,z,r）为关节坐标系下的目标点 坐标</p>
<p>//MOVJ 模式，（x,y,z,r）为关节坐标系下的目标点 坐标</p>
<p>//MOVL 模式，（x,y,z,r）为关节坐标系下的目标</p>
<p>点坐标</p>
<p>//MOVJ 模式，（x,y,z,r）为关节坐标系下的坐标</p>
<p>增量</p>
<p>//MOVL 模式，（x,y,z,r）为笛卡尔坐标系下的坐</p>
<p>标增量</p>
</blockquote></td>
</tr>
<tr class="odd">
<td colspan="2"><blockquote>
<p>MOVJ_XYZ_INC, //MOVJ 模式，（x,y,z,r）为笛卡尔坐标系下的坐 标增量</p>
<p>JUMP_MOVL_XYZ, //JUMP 模式，平移时运动模式为 MOVL。</p>
<p>（x,y,z,r）为笛卡尔坐标系下的坐标增量</p>
<p>};</p>
<p>ptpCmd：PTPCmd 指针</p>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark83
> .anchor}![](media/image255.png){width="0.40999890638670167in"
> height="0.12666557305336834in"} **执行带** **I/O 控制的** **PTP 指令**
>
> 表 12.14 执行带 I/O 控制的 PTP 指令

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetPTPPOCmd(PTPCmd *ptpCmd, ParallelOutputCmd *parallelCmd, int
parallelCmdCount, bool isQueued,uint64_t *queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>执行带 I/O 控制的 PTP 指令。I/O 说明请参见《Dobot Magician
用户手册》</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 27%" />
<col style="width: 2%" />
<col style="width: 60%" />
</colgroup>
<tbody>
<tr class="odd">
<td rowspan="4"><blockquote>
<p>参数</p>
</blockquote></td>
<td colspan="3"><blockquote>
<p>PTPCmd 定义：</p>
<p>typedef struct tagPTPCmd {</p>
<p>uint8_tptpMode; //PTP 模式，取值范围：0~9</p>
<p>float x; //（x,y,z,r）为坐标参数，可为笛卡尔坐标、关节坐</p>
<p>标、 笛卡尔坐标增量或关节坐标增量</p>
<p>float y;</p>
<p>float z;</p>
<p>float r;</p>
<p>}PTPCmd;</p>
<p>//其中， ptpMode 取值如下：</p>
<p>enum {</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>JUMP_XYZ,</p>
<p>MOVJ_XYZ,</p>
<p>MOVL_XYZ,</p>
<p>JUMP_ANGLE,</p>
<p>MOVJ_ANGLE,</p>
<p>MOVL_ANGLE,</p>
<p>MOVJ_INC,</p>
<p>MOVL_INC,</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>//JUMP 模式，（x,y,z,r）为笛卡尔坐标系下的目标 点坐标</p>
<p>//MOVJ 模式，（x,y,z,r）为笛卡尔坐标系下的目标 点坐标</p>
<p>//MOVL 模式，（x,y,z,r）为笛卡尔坐标系下的目标 点坐标</p>
<p>//JUMP 模式，（x,y,z,r）为关节坐标系下的目标点 坐标</p>
<p>//MOVJ 模式，（x,y,z,r）为关节坐标系下的目标点 坐标</p>
<p>//MOVL 模式，（x,y,z,r）为关节坐标系下的目标</p>
<p>点坐标</p>
<p>//MOVJ 模式，（x,y,z,r）为关节坐标系下的坐标</p>
<p>增量</p>
<p>//MOVL 模式，（x,y,z,r）为笛卡尔坐标系下的坐</p>
<p>标增量</p>
</blockquote></td>
</tr>
<tr class="odd">
<td colspan="3"><blockquote>
<p>MOVJ_XYZ_INC, //MOVJ 模式，（x,y,z,r）为笛卡尔坐标系下的坐 标增量</p>
<p>JUMP_MOVL_XYZ, //JUMP 模式，平移时运动模式为 MOVL。</p>
<p>（x,y,z,r）为笛卡尔坐标系下的坐标增量</p>
<p>};</p>
<p>ParallelOutputCmd 定义：</p>
<p>typedef struct tagParallelOutputCmd {</p>
</blockquote></td>
</tr>
<tr class="even">
<td colspan="2"><blockquote>
<p>uint8_t ratio;</p>
<p>uint16_t address;</p>
<p>uint8_t level;</p>
</blockquote></td>
<td><blockquote>
<p>//设置运动时两点之间距离的百分比，即在该 位置触发 I/O</p>
<p>//I/O 地址。取值范围：1~20</p>
<p>//输出值</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td></td>
<td><blockquote>
<p>}ParallelOutputCmd;</p>
<p>ptpCmd：PTPCmd 指针</p>
<p>parallelCmd：ParallelOutputCmd 指针</p>
<p>parallelCmdCount：I/O 个数</p>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark84
> .anchor}![](media/image256.png){width="0.4083333333333333in"
> height="0.12999890638670167in"} **执行带滑轨的** **PTP 指令**
>
> 表 12.15 执行带滑轨的 PTP 指令

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 26%" />
<col style="width: 64%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>int SetPTPWithLCmd(PTPWithLCmd *ptpWithLCmd, bool isQueued,uint64_t
*queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>执行带滑轨的 PTP 指令</p>
</blockquote></td>
</tr>
<tr class="odd">
<td rowspan="4"><blockquote>
<p>参数</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>PTPWithLCmd 定义：</p>
<p>typedef struct tagPTPWithLCmd {</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>uint8_tptpMode;</p>
<p>float x;</p>
<p>float y;</p>
<p>float z;</p>
<p>float r;</p>
<p>float l;</p>
</blockquote></td>
<td><blockquote>
<p>//PTP 模式，取值范围：0~9</p>
<p>//（x,y,z,r）为坐标参数，可为笛卡尔坐标、关节坐 标、
笛卡尔坐标增量或关节坐标增量</p>
<p>//滑轨运行距离</p>
</blockquote></td>
</tr>
<tr class="odd">
<td colspan="2"><blockquote>
<p>}PTPWithLCmd;</p>
<p>//其中， ptpMode 取值如下：</p>
<p>enum {</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>JUMP_XYZ,</p>
<p>MOVJ_XYZ,</p>
<p>MOVL_XYZ,</p>
<p>JUMP_ANGLE,</p>
</blockquote></td>
<td><blockquote>
<p>//JUMP 模式，（x,y,z,r）为笛卡尔坐标系下的目标 点坐标</p>
<p>//MOVJ 模式，（x,y,z,r）为笛卡尔坐标系下的目标 点坐标</p>
<p>//MOVL 模式，（x,y,z,r）为笛卡尔坐标系下的目标 点坐标</p>
<p>//JUMP 模式，（x,y,z,r）为关节坐标系下的目标点</p>
<p>坐标</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td></td>
<td><blockquote>
<p>MOVJ_ANGLE, //MOVJ 模式，（x,y,z,r）为关节坐标系下的目标点</p>
<p>坐标</p>
<p>MOVL_ANGLE, //MOVL 模式，（x,y,z,r）为关节坐标系下的目标</p>
<p>点坐标</p>
<p>MOVJ_INC, //MOVJ 模式，（x,y,z,r）为关节坐标系下的坐标</p>
<p>增量</p>
<p>MOVL_INC, //MOVL 模式，（x,y,z,r）为笛卡尔坐标系下的坐</p>
<p>标增量</p>
<p>MOVJ_XYZ_INC, //MOVJ 模式，（x,y,z,r）为笛卡尔坐标系下的坐 标增量</p>
<p>JUMP_MOVL_XYZ, //JUMP 模式，平移时运动模式为 MOVL。</p>
<p>（x,y,z,r）为笛卡尔坐标系下的坐标增量</p>
<p>};</p>
<p>ptpWithLCmd：PTPWithLCmd 指针</p>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark85
> .anchor}![](media/image257.png){width="0.4083333333333333in"
> height="0.13in"} **执行带** **I/O 控制和滑轨的** **PTP 指令**
>
> 表 12.16 执行带 I/O 控制和滑轨的 PTP 指令

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 27%" />
<col style="width: 63%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>int SetPTPPOWithLCmd(PTPWithLCmd *ptpWithLCmd, ParallelOutputCmd</p>
<p>*parallelCmd,int parallelCmdCount, bool isQueued,uint64_t
*queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>执行带I/O控制和滑轨的PTP指令</p>
</blockquote></td>
</tr>
<tr class="odd">
<td rowspan="3"><blockquote>
<p>参数</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>PTPWithLCmd 定义：</p>
<p>typedef struct tagPTPWithLCmd {</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>uint8_tptpMode;</p>
<p>float x;</p>
<p>float y;</p>
<p>float z;</p>
<p>float r;</p>
<p>float l;</p>
</blockquote></td>
<td><blockquote>
<p>//PTP 模式，取值范围：0~9</p>
<p>//（x,y,z,r）为坐标参数，可为笛卡尔坐标、关节坐 标、
笛卡尔坐标增量或关节坐标增量</p>
<p>//滑轨运行距离</p>
</blockquote></td>
</tr>
<tr class="odd">
<td colspan="2"><blockquote>
<p>}PTPWithLCmd;</p>
<p>//其中， ptpMode 取值如下：</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 28%" />
<col style="width: 62%" />
</colgroup>
<tbody>
<tr class="odd">
<td rowspan="2"></td>
<td><blockquote>
<p>enum {</p>
<p>JUMP_XYZ,</p>
<p>MOVJ_XYZ,</p>
<p>MOVL_XYZ,</p>
<p>JUMP_ANGLE,</p>
<p>MOVJ_ANGLE,</p>
<p>MOVL_ANGLE,</p>
<p>MOVJ_INC,</p>
<p>MOVL_INC,</p>
<p>MOVJ_XYZ_INC,</p>
</blockquote></td>
<td><blockquote>
<p>//JUMP 模式，（x,y,z,r）为笛卡尔坐标系下的目标 点坐标</p>
<p>//MOVJ 模式，（x,y,z,r）为笛卡尔坐标系下的目标 点坐标</p>
<p>//MOVL 模式，（x,y,z,r）为笛卡尔坐标系下的目标 点坐标</p>
<p>//JUMP 模式，（x,y,z,r）为关节坐标系下的目标点</p>
<p>坐标</p>
<p>//MOVJ 模式，（x,y,z,r）为关节坐标系下的目标点</p>
<p>坐标</p>
<p>//MOVL 模式，（x,y,z,r）为关节坐标系下的目标 点坐标</p>
<p>//MOVJ 模式，（x,y,z,r）为关节坐标系下的坐标 增量</p>
<p>//MOVL 模式，（x,y,z,r）为笛卡尔坐标系下的坐 标增量</p>
<p>//MOVJ 模式，（x,y,z,r）为笛卡尔坐标系下的坐 标增量</p>
</blockquote></td>
</tr>
<tr class="even">
<td colspan="2"><blockquote>
<p>JUMP_MOVL_XYZ, //JUMP 模式，平移时运动模式为 MOVL。</p>
<p>（x,y,z,r）为笛卡尔坐标系下的坐标增量</p>
<p>};</p>
<p>ParallelOutputCmd 定义：</p>
<p>typedef struct tagParallelOutputCmd {</p>
</blockquote>
<table>
<colgroup>
<col style="width: 29%" />
<col style="width: 70%" />
</colgroup>
<tbody>
<tr class="odd">
<td>uint8_t ratio;</td>
<td><blockquote>
<p>//设置运动时两点之间距离的百分比，即在该 位置触发 I/O</p>
</blockquote></td>
</tr>
<tr class="even">
<td>uint16_t address;</td>
<td><blockquote>
<p>//I/O 地址。取值范围：1~20</p>
</blockquote></td>
</tr>
</tbody>
</table>
<blockquote>
<p>uint8_t level; //输出值</p>
<p>}ParallelOutputCmd;</p>
<p>ptpCmd：PTPCmd 指针</p>
<p>parallelCmd：ParallelOutputCmd 指针</p>
<p>parallelCmdCount：I/O 个数</p>
<p>isQueued:是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark86
> .anchor}![](media/image259.png){width="0.21666557305336834in"
> height="0.14333333333333334in"} **CP 功能**
>
> CP即连续运动轨迹。
>
> []{#_bookmark87
> .anchor}![](media/image260.png){width="0.2949989063867017in"
> height="0.12999890638670167in"} **设置** **CP 运动的速度和加速度**
>
> 表 13.1 设置连续轨迹功能参数接口说明

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetCPParams(CPParams *cpParams, bool isQueued,uint64_t</p>
<p>*queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置CP运动的速度（mm/s）和加速度（mm/s2）</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>CPParams 定义：</p>
<p>typedef struct tagCPParams {</p>
</blockquote>
<table>
<colgroup>
<col style="width: 31%" />
<col style="width: 68%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>float planAcc;</p>
</blockquote></td>
<td><blockquote>
<p>//规划加速度最大值</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>float junctionVel;</p>
</blockquote></td>
<td><blockquote>
<p>//拐角速度最大值</p>
</blockquote></td>
</tr>
<tr class="odd">
<td>union {</td>
<td></td>
</tr>
<tr class="even">
<td><blockquote>
<p>float acc;</p>
</blockquote></td>
<td>//实际加速度最大值，非实时模式时有效</td>
</tr>
<tr class="odd">
<td><blockquote>
<p>float period;</p>
</blockquote></td>
<td><blockquote>
<p>//插补周期，实时模式时有效</p>
</blockquote></td>
</tr>
</tbody>
</table>
<blockquote>
<p>};</p>
<p>uint8_trealTimeTrack; //0：非实时模式：所有指令下发后再运行</p>
<p>1：实时模式：边下发指令边运行</p>
<p>}CPParams;</p>
<p>cpParams：CPParams 指针</p>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回，导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark88
> .anchor}![](media/image261.png){width="0.313332239720035in"
> height="0.13in"} **获取** **CP 运动的速度和加速度**
>
> 表 13.2 获取连续轨迹功能参数接口说明

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 35%" />
<col style="width: 55%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>int GetCPParams(CPParams *cpParams)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>获取连续轨迹模型下相关参数</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>CPParams 定义：</p>
<p>typedef struct tagCPParams {</p>
<p>float planAcc;</p>
<p>float junctionVel;</p>
</blockquote></td>
<td><blockquote>
<p>//规划加速度最大值</p>
<p>//拐角速度最大值</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 33%" />
<col style="width: 56%" />
</colgroup>
<tbody>
<tr class="odd">
<td></td>
<td><blockquote>
<p>union {</p>
<p>float acc;</p>
<p>float period;</p>
<p>};</p>
<p>uint8_trealTimeTrack;</p>
<p>}CPParams;</p>
<p>cpParams：CPParams 指针</p>
</blockquote></td>
<td><blockquote>
<p>//实际加速度最大值，非实时模式时有效 //插补周期，实时模式时有效</p>
<p>//0：非实时模式：所有指令下发后再运行</p>
<p>1：实时模式：边下发指令边运行</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark89 .anchor}![](media/image262.png){width="0.315in"
> height="0.12999890638670167in"} **执行** **CP 指令**
>
> 表 13.3 执行 CP 指令

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 31%" />
<col style="width: 59%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>int SetCPCmd(CPCmd *cpCmd, bool isQueued, uint64_t
*queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>执行 CP 指令</p>
</blockquote></td>
</tr>
<tr class="odd">
<td rowspan="2"><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>CPCmd 定义：</p>
<p>typedef struct tagCPCmd {</p>
<p>uint8_t cpMode;</p>
</blockquote></td>
<td><blockquote>
<p>//CP 模式。0：相对模式，表示相对距离，即笛
卡尔坐标增量。1：绝对模式，表示绝对距离，</p>
</blockquote></td>
</tr>
<tr class="even">
<td colspan="2"><blockquote>
<p>即笛卡尔坐标系下目标点坐标</p>
<p>float x; //x,y,z 可以设置为坐标增量，也可设置为目的坐</p>
<p>标点</p>
<p>float y;</p>
<p>float z;</p>
<p>union {</p>
<p>float velocity; //保留</p>
<p>float power; //保留</p>
<p>};</p>
<p>}CPCmd;</p>
<p>cpCmd：CPCmd 指针</p>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> ![](media/image264.png){width="0.2772594050743657in"
> height="0.2393788276465442in"}注意
>
> 当指令队列中有多条连续的CP指令时，
> Dobot控制器将自动前瞻。前瞻的条件是， 队列中这些CP指令之间没有JOG
> 、PTP 、ARC 、WAIT 、TRIG等指令。
>
> []{#_bookmark90
> .anchor}![](media/image265.png){width="0.318332239720035in"
> height="0.12999890638670167in"} **执行带激光雕刻的** **CP 指令**
>
> 表 13.4 执行带激光雕刻的 CP 指令

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetCPLECmd (CPCmd *cpCmd, bool isQueued,uint64_t
*queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>执行带激光雕刻的 CP 指令。此时激光功率将在运动命令下发时起效</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>CPCmd 定义：</p>
<p>typedef struct tagCPCmd {</p>
<p>uint8_t cpMode; //CP 模式。0：相对模式，表示相对距离，即笛</p>
<p>卡尔坐标增量。 1：绝对模式，表示绝对距离，
即笛卡尔坐标系下目标点坐标</p>
<p>float x; //x,y,z 可以设置为坐标增量，也可设置为目的坐</p>
<p>标点</p>
<p>float y;</p>
<p>float z;</p>
<p>union {</p>
<p>float velocity; //保留</p>
<p>float power; //激光功率：0~100</p>
<p>}</p>
<p>}CPCmd;</p>
<p>cpCmd：CPCmd 指针</p>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark91
> .anchor}![](media/image267.png){width="0.21666557305336834in"
> height="0.14333333333333334in"} **ARC 功能**
>
> 圆弧模式即示教后存点回放的运动轨迹为圆弧。圆弧轨迹是空间的圆弧，由当前点、圆
> 弧上任一点和圆弧结束点三点共同确定。圆弧总是从起点经过圆弧上任一点再到结束点，如
> [图 14.1所](#_bookmark92)示。
>
> ![](media/image268.png){width="0.2772594050743657in"
> height="0.23937992125984253in"}注意
>
> 使用圆弧运动模式时，需结合其他运动模式确认圆弧上的三点，且三点不能在同
> 一条直线上。

![](media/image269.jpeg){width="4.3533333333333335in"
height="1.8149989063867016in"}

> 图 14.1 圆弧运动模式
>
> []{#_bookmark92
> .anchor}![](media/image270.png){width="0.2949989063867017in"
> height="0.12666557305336834in"} **设置** **ARC 运动的速度和加速度**
>
> 表 14.1 设置 ARC 运动的速度和加速度

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetARCParams(ARCParams *arcParams, bool isQueued,uint64_t
*queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置 ARC 运动的速度（mm/s）和加速度（mm/s2）</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>ARCParams 定义：</p>
<p>typedef struct tagARCParams {</p>
<p>float xyzVelocity; //PTP 模式下 X,Y,Z 3 轴坐标轴速度</p>
<p>float rVelocity; //PTP 模式下末端 R 轴速度</p>
<p>float xyzAcceleration; //PTP 模式下 X,Y,Z 3 轴坐标轴加速度</p>
<p>float rAccleration; //PTP 模式下末端 R 轴加速度</p>
<p>} ARCParams;</p>
<p>arcParams：圆弧插补功能参数</p>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td></td>
<td><blockquote>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark93
> .anchor}![](media/image271.png){width="0.313332239720035in"
> height="0.12666447944007in"} **获取** **ARC 运动的速度和加速度**
>
> 表 14.2 获取 ARC 运动的速度和加速度

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetARCParams(ARCParams *arcParams)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取 ARC 运动的速度（mm/s）和加速度（mm/s2）</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>ARCParams 定义：</p>
<p>typedef struct tagARCParams {</p>
<p>float xyzVelocity; //PTP 模式下 X,Y,Z 3 轴坐标轴速度</p>
<p>float rVelocity; //PTP 模式下末端 R 轴速度</p>
<p>float xyzAcceleration; //PTP 模式下 X,Y,Z 3 轴坐标轴加速度</p>
<p>float rAccleration; //PTP 模式下末端 R 轴加速度</p>
<p>} ARCParams;</p>
<p>arcParams：ARCParams 指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark94 .anchor}![](media/image272.png){width="0.315in"
> height="0.13in"} **执行** **ARC 指令**
>
> 表 14.3 执行圆弧插补功能接口说明

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetARCCmd(ARCCmd *arcCmd, bool isQueued,uint64_t</p>
<p>*queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>执行 ARC 指令。设置 ARC 运动的速度和加速度后，调用该函数可使机械 臂按
ARC 模式运动至目标点</p>
<p>该运动模式需结合其他模式一起使用，构成圆弧上三点</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>ARCCmd 定义：</p>
<p>typedef struct tagARCCmd {</p>
<p>struct {</p>
<p>float x;</p>
<p>float y;</p>
<p>float z;</p>
<p>float r;</p>
<p>}cirPoint; //圆弧中间点，需设置为笛卡尔坐标</p>
<p>struct {</p>
<p>float x;</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 17%" />
<col style="width: 73%" />
</colgroup>
<tbody>
<tr class="odd">
<td rowspan="2"></td>
<td><blockquote>
<p>float float float</p>
<p>}toPoint;</p>
</blockquote></td>
<td><blockquote>
<p>y;</p>
<p>z;</p>
<p>r;</p>
<p>//圆弧目标点，需设置为笛卡尔坐标</p>
</blockquote></td>
</tr>
<tr class="even">
<td colspan="2"><blockquote>
<p>}ARCCmd;</p>
<p>arcCmd：ARCCmd 指针</p>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark95
> .anchor}![](media/image273.png){width="0.318332239720035in"
> height="0.12666557305336834in"} **执行** **CIRCLE 指令**
>
> 圆形模式与圆弧模式相似，示教后存点回放的运动轨迹为整圆。使用圆形模式时，也需
> 结合其他运动模式确认圆形上的三点。
>
> 表 14.4 执行 CIRCLE 指令

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetCircleCmd(CircleCmd *circleCmd, bool isQueued,uint64_t
*queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>执行 CIRCLE 指令。设置整圆运动的速度和加速度后，调用该函数可使机械
臂按 CIRCLE 模式运动至目标点</p>
<p>该运动模式需结合其他模式一起使用，构成圆上三点</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>CircleCmd 定义：</p>
<p>typedef struct tagCircleCmd{</p>
<p>struct {</p>
<p>float x;</p>
<p>float y;</p>
<p>float z;</p>
<p>float r;</p>
<p>}cirPoint; //圆弧中间点，需设置为笛卡尔坐标</p>
<p>struct {</p>
<p>float x;</p>
<p>float y;</p>
<p>float z;</p>
<p>float r;</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 26%" />
<col style="width: 63%" />
</colgroup>
<tbody>
<tr class="odd">
<td rowspan="2"></td>
<td><blockquote>
<p>}toPoint;</p>
<p>uint32_t count;</p>
</blockquote></td>
<td><blockquote>
<p>//圆弧目标点，需设置为笛卡尔坐标</p>
<p>//整圆个数</p>
</blockquote></td>
</tr>
<tr class="even">
<td colspan="2"><blockquote>
<p>}CircleCmd;</p>
<p>circleCmd：CircleCmd 指针</p>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark96
> .anchor}![](media/image275.png){width="0.21666557305336834in"
> height="0.14333333333333334in"} **丢步检测功能**
>
> []{#_bookmark97
> .anchor}![](media/image276.png){width="0.2949989063867017in"
> height="0.12999890638670167in"} **设置丢步检测阈值**
>
> 表 15.1 设置丢步检测阈值接口说明

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetLostStepParams(float threshold, bool isQueued,uint64_t</p>
<p>*queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置丢步检测阈值，用于检测定位误差是否超过该阈值。如果超过该阈值，
则说明电机丢步</p>
<p>如果用户不调用该接口， 则丢步检测阈值默认为 5</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>threshold：阈值</p>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark98
> .anchor}![](media/image277.png){width="0.313332239720035in"
> height="0.13in"} **执行丢步检测**
>
> 表 15.2 设置丢步检测阈值接口说明

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetLostStepCmd(bool isQueued,uint64_t *queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>执行丢步检测， 当检测到丢步时会停止命令队列。该指令只能作为队列指
令， 即“isQueued ”必须设置为“ 1 ”</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark99 .anchor}![](media/image278.png){width="0.315in"
> height="0.13in"} **示例：** **丢步检测**
>
> 程序 15.1 丢步检测
>
> int main(void)
>
> {
>
> uint64_t queuedCmdIndex = 0;
>
> PTPCmd cmd;
>
> cmd.ptpMode = 0;
>
> cmd.x
>
> cmd.y
>
> cmd.z
>
> cmd.r
>
> = 200;
>
> = 0;
>
> = 0;

= 0;

> ConnectDobot(NULL, 115200, NULL, NULL, NULL);
>
> SetQueuedCmdStartExec();
>
> SetPTPCmd(&cmd, true, &queuedCmdIndex);
>
> SetLostStepCmd(true, &queuedCmdIndex)
>
> SetQueuedCmdStopExec();
>
> DisconnectDobot();
>
> }
>
> []{#_bookmark100
> .anchor}![](media/image280.png){width="0.21666557305336834in"
> height="0.14333333333333334in"} **WAIT 功能**
>
> []{#_bookmark101
> .anchor}![](media/image281.png){width="0.2949989063867017in"
> height="0.12999890638670167in"} **执行时间等待指令**
>
> 表 16.1 执行时间等待功能接口说明

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetWAITCmd(WAITCmd *waitCmd, bool isQueued,uint64_t</p>
<p>*queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>执行时间等待指令。如果需设置前个指令运行后的暂停时间， 可调用此函数
该指令只能作为队列指令，“isQueued ”必须设置为“ 1 ”。将该命令设置
为立即指令可能会导致正在执行的 WAIT 队列指令的暂停时间变化</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>WAITCmd 定义：</p>
<p>typedef struct tagWAITCmd {</p>
<p>uint32_ttimeout; //单位：ms</p>
<p>}WAITCmd;</p>
<p>waitCmd：时间等待功能变量</p>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark102
> .anchor}![](media/image283.png){width="0.21666557305336834in"
> height="0.14333333333333334in"} **执行触发指令**
>
> 表 17.1 执行触发指令

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 28%" />
<col style="width: 62%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>int SetTRIGCmd(TRIGCmd *trigCmd, bool isQueued,uint64_t</p>
<p>*queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>执行触发指令</p>
<p>该指令只能作为队列指令，“isQueued ”必须设置为“ 1 ”。将该命令设置
为立即指令可能会导致正在执行的 TRIG 队列指令的触发条件变化。</p>
</blockquote></td>
</tr>
<tr class="odd">
<td rowspan="3"><blockquote>
<p>参数</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>TRIGCmd 定义：</p>
<p>typedef struct tagTRIGCmd {</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>uint8_t address;</p>
<p>uint8_t mode;</p>
<p>uint8_t condition;</p>
<p>uint16_t threshold;</p>
</blockquote></td>
<td><blockquote>
<p>//I/O 地址。取值范围： 1~20</p>
<p>//触发模式。0：电平触发。1：A/D 触发</p>
<p>//触发条件。</p>
<p>电平：0 ，等于。1 ，不等于</p>
<p>A/D：0 ，小于。1 ，小于等于</p>
<p>2 ，大于等于。3 ，大于</p>
<p>//触发阈值： 电平，0 或 1 。A/D：0~4095</p>
</blockquote></td>
</tr>
<tr class="odd">
<td colspan="2"><blockquote>
<p>}TRIGCmd;</p>
<p>trigCmd：TRIGCmd 指针</p>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark103
> .anchor}![](media/image285.png){width="0.21666557305336834in"
> height="0.14333333333333334in"} **EIO 功能**
>
> 在 Dobot 控制器中，所有的扩展 I/O 都是统一编址的。根据现有情况， I/O
> 的功能包括 以下内容：
>
> . 高低电平输出功能。
>
> . PWM输出功能。
>
> . 读取输入高低电平功能。
>
> . 读取输入模数转换值功能。
>
> 部分 I/O 可能同时具有以上的功能。在使用不同的功能时，需要先配置 I/O
> 的复用。I/O 详细说明， 请参见《Dobot Magician 用户手册》。
>
> []{#_bookmark104
> .anchor}![](media/image286.png){width="0.2949989063867017in"
> height="0.12999890638670167in"} **设置** **I/O 复用**
>
> 表 18.1 设置 I/O 复用

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 31%" />
<col style="width: 59%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>int SetIOMultiplexing(IOMultiplexing *ioMultiplexing, bool
isQueued,uint64_t *queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>设置 I/O 复用。在使用 I/O 前需要设置 I/O 复用功能。</p>
</blockquote></td>
</tr>
<tr class="odd">
<td rowspan="3"><blockquote>
<p>参数</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>IOMultiplexing 定义：</p>
<p>typedef struct tagIOMultiplexing {</p>
<p>uint8_t address; //I/O 地址。取值范围： 1~20</p>
<p>uint8_t multiplex; //IO 功能。取值范围：0~6</p>
<p>}IOMultiplexing;</p>
<p>其中 mutiplex 支持的取值如下所示：</p>
<p>typedefenum tagIOFunction {</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>IOFunctionDummy,</p>
<p>IOFunctionDO,</p>
<p>IOFunctionPWM,</p>
<p>IOFunctionDI,</p>
<p>IOFunctionADC IOFunctionDIPU</p>
<p>IOFunctionDIPD</p>
</blockquote></td>
<td><blockquote>
<p>//不配置功能</p>
<p>//I/O 输出</p>
<p>//PWM 输出</p>
<p>//I/O 输入</p>
<p>//A/D 输入</p>
<p>//上拉输入</p>
<p>//下拉输入</p>
</blockquote></td>
</tr>
<tr class="odd">
<td colspan="2"><blockquote>
<p>} IOFunction;</p>
<p>ioMultiplexing：IOMultiplexing 指针</p>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td></td>
<td><blockquote>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark105
> .anchor}![](media/image287.png){width="0.313332239720035in"
> height="0.12999890638670167in"} **读取** **I/O 复用**
>
> 表 18.2 读取 I/O 复用

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 32%" />
<col style="width: 58%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>int GetIOMultiplexing(IOMultiplexing *ioMultiplexing)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>读取 I/O 复用</p>
</blockquote></td>
</tr>
<tr class="odd">
<td rowspan="3"><blockquote>
<p>参数</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>IOMultiplexing 定义：</p>
<p>typedef struct tagIOMultiplexing {</p>
<p>uint8_t address; //I/O 地址</p>
<p>uint8_t multiplex; //IO 功能。取值范围：0~6</p>
<p>}IOMultiplexing;</p>
<p>其中 mutiplex 支持的取值如下所示：</p>
<p>typedefenum tagIOFunction {</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>IOFunctionDummy,</p>
<p>IOFunctionDO,</p>
<p>IOFunctionPWM,</p>
<p>IOFunctionDI,</p>
<p>IOFunctionADC IOFunctionDIPU</p>
<p>IOFunctionDIPD</p>
</blockquote></td>
<td><blockquote>
<p>//不配置功能</p>
<p>//I/O 输出</p>
<p>//PWM 输出</p>
<p>//I/O 输入</p>
<p>//A/D 输入</p>
<p>//上拉输入</p>
<p>//下拉输入</p>
</blockquote></td>
</tr>
<tr class="odd">
<td colspan="2"><blockquote>
<p>} IOFunction;</p>
<p>ioMultiplexing：IOMultiplexing 指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark106 .anchor}![](media/image288.png){width="0.315in"
> height="0.12999890638670167in"} **设置** **I/O 输出电平**
>
> 表 18.3 设置 I/O 输出电平

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 26%" />
<col style="width: 63%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>int SetIODO(IODO *ioDO, bool isQueued,uint64_t *queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>设置 I/O 输出电平</p>
</blockquote></td>
</tr>
<tr class="odd">
<td rowspan="2"><blockquote>
<p>参数</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>IODO 定义：</p>
<p>typedef struct tagIODO {</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>uint8_t address;</p>
<p>uint8_t level;</p>
<p>}IODO;</p>
<p>ioDO：IODO 指针</p>
</blockquote></td>
<td><blockquote>
<p>//I/O 地址</p>
<p>//输出电平。0：低电平。1：高电平</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td></td>
<td><blockquote>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark107
> .anchor}![](media/image289.png){width="0.318332239720035in"
> height="0.12999890638670167in"} **读取** **I/O 输出电平**
>
> 表 18.4 读取 I/O 输出电平

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 26%" />
<col style="width: 63%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>int GetIODO(IODO *ioDO)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>读取 I/O 输出电平</p>
</blockquote></td>
</tr>
<tr class="odd">
<td rowspan="2"><blockquote>
<p>参数</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>IODO 定义：</p>
<p>typedef struct tagIODO {</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>uint8_t address;</p>
<p>uint8_t level;</p>
<p>}IODO;</p>
<p>ioDO：IODO 指针</p>
</blockquote></td>
<td><blockquote>
<p>//I/O 地址</p>
<p>//输出电平。0：低电平。1：高电平</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark108
> .anchor}![](media/image290.png){width="0.31666557305336834in"
> height="0.13in"} **设置** **PWM 输出**
>
> 表 18.5 设置 PWM 输出

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetIOPWM(IOPWM *ioPWM, bool isQueued,uint64_t
*queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置 PWM 输出</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>IOPWM 定义：</p>
<p>typedef struct tagIOPWM {</p>
<p>uint8_t address; //I/O 地址</p>
<p>float frequency; //PWM 频率。取值范围： 10HZ~1MHz</p>
<p>float dutyCycle; //PWM 占空比。取值范围：0~100</p>
<p>}IOPWM;</p>
<p>ioPWM：IOPWM 指针</p>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark109
> .anchor}![](media/image291.png){width="0.31666557305336834in"
> height="0.12999890638670167in"} **读取** **PWM 输出**
>
> 表 18.6 读取 PWM 输出

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 28%" />
<col style="width: 61%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>int GetIOPWM(IOPWM *ioPWM)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>读取 PWM 输出</p>
</blockquote></td>
</tr>
<tr class="odd">
<td rowspan="2"><blockquote>
<p>参数</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>IOPWM 定义：</p>
<p>typedef struct tagIOPWM {</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>uint8_t address;</p>
<p>float frequency;</p>
<p>float dutyCycle;</p>
<p>}IOPWM;</p>
<p>ioPWM：IOPWM 指针</p>
</blockquote></td>
<td><blockquote>
<p>//I/O 地址</p>
<p>//PWM 频率。取值范围： 10HZ~1MHz</p>
<p>//PWM 占空比。取值范围：0~100</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark110 .anchor}![](media/image292.png){width="0.315in"
> height="0.12999890638670167in"} **读取** **I/O 输入电平**
>
> 表 18.7 读取 I/O 输入电平

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 28%" />
<col style="width: 61%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>int GetIODI(IODI *ioDI)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>读取 I/O 输入电平</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>IODI 定义：</p>
<p>typedef struct tagIODI {</p>
<p>uint8_t address;</p>
<p>uint8_t level;</p>
<p>}IODI;</p>
<p>ioDI：IODI 指针</p>
</blockquote></td>
<td><blockquote>
<p>//I/O 地址</p>
<p>//输入电平。0：低电平。1：高电平</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark111 .anchor}![](media/image293.png){width="0.315in"
> height="0.13in"} **读取** **A/D 输入**
>
> 表 18.8 读取 A/D 输入

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 28%" />
<col style="width: 61%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>int GetIOADC(IOADC *ioADC)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>读取 A/D 输入</p>
</blockquote></td>
</tr>
<tr class="odd">
<td rowspan="2"><blockquote>
<p>参数</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>IOADC 定义：</p>
<p>typedef struct tagIOADC {</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>uint8_t address;</p>
<p>uint16_t value;</p>
<p>}IOADC;</p>
<p>ioADC：IOADC 指针</p>
</blockquote></td>
<td><blockquote>
<p>//I/O 地址</p>
<p>//A/D 输入值。取值范围：0~4095</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td colspan="2"><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark112
> .anchor}![](media/image294.png){width="0.313332239720035in"
> height="0.12999890638670167in"} **设置扩展电机速度**
>
> 表 18.9 设置扩展电机速度

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetEMotor(EMotor *eMotor, bool isQueued,uint64_t
*queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置扩展电机速度。调用此函数后电机会以一定的速度不停的运行</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>EMotor 定义：</p>
<p>typedef struct tagEMotor {</p>
</blockquote>
<table>
<colgroup>
<col style="width: 31%" />
<col style="width: 68%" />
</colgroup>
<tbody>
<tr class="odd">
<td>uint8_t index;</td>
<td>//电机编号。0：Stepper1 。 1：Stepper2</td>
</tr>
<tr class="even">
<td>uint8_tisEnabled;</td>
<td><blockquote>
<p>//电机控制使能。0：未使能。 1：使能</p>
</blockquote></td>
</tr>
<tr class="odd">
<td>uint32_t speed;</td>
<td><blockquote>
<p>//电机控制速度（脉冲个数每秒）</p>
</blockquote></td>
</tr>
</tbody>
</table>
<blockquote>
<p>}EMotor;</p>
<p>eMotor：EMoto 指针</p>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark113
> .anchor}![](media/image295.png){width="0.40499890638670166in"
> height="0.13in"} **设置扩展电机速度和移动距离**
>
> 表 18.10 设置扩展电机速度和移动距离接口说明

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetEMotorS(EMotorS *eMotorS, bool isQueued,uint64_t
*queuedCmdIndex)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置扩展电机速度和移动距离。当需要以一定速度运行一段距离时可调用此
函数</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>EMotorS 定义：</p>
<p>typedef struct tagEMotorS {</p>
</blockquote>
<table>
<colgroup>
<col style="width: 31%" />
<col style="width: 68%" />
</colgroup>
<tbody>
<tr class="odd">
<td>uint8_t index;</td>
<td>//电机编号。0：Stepper1 。 1：Stepper2</td>
</tr>
<tr class="even">
<td>uint8_tisEnabled;</td>
<td><blockquote>
<p>//电机控制使能。0：未使能。 1：使能</p>
</blockquote></td>
</tr>
<tr class="odd">
<td>uint32_t speed;</td>
<td><blockquote>
<p>//电机控制速度（脉冲个数每秒）</p>
</blockquote></td>
</tr>
<tr class="even">
<td>uint32_t distance;</td>
<td><blockquote>
<p>//电机移动距离(脉冲个数)</p>
</blockquote></td>
</tr>
</tbody>
</table>
<blockquote>
<p>}EMotorS;</p>
<p>eMotorS：EMotorS 指针</p>
<p>isQueued：是否将该指令指定为队列命令</p>
<p>queuedCmdIndex：若选择将指令加入队列，则表示指令在队列的索引号。
否则， 该参数无意义</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_BufferFull：指令队列已满</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark114
> .anchor}![](media/image296.png){width="0.38666666666666666in"
> height="0.12999890638670167in"} **使能光电传感器**
>
> 表 18.11 使能光电传感器

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetInfaredSensor(bool enable,InfraredPort infraredPort, uint8_t
version)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>使能光电传感器</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>InfraredPort 定义：</p>
<p>enum InfraredPort{</p>
<p>IF_PORT_GP1;</p>
<p>IF_PORT_GP2;</p>
<p>IF_PORT_GP4;</p>
<p>IF_PORT_GP5;</p>
<p>};</p>
<p>enable：使能标志。0：未使能。 1：使能</p>
<p>infraredPort：光电传感器连接至机械臂的接口，请选择对应的接口</p>
<p>version：光电传感器版本号。0：版本为 V1.0 。1：版本为 V2.0</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark115
> .anchor}![](media/image297.png){width="0.40499890638670166in"
> height="0.13in"} **获取光电传感器读数**
>
> 表 18.12 获取光电传感器读数

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetInfraredSensor (InfraredPort infraredPort, uint8_t *value)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取光电传感器读数</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>InfraredPort 定义：</p>
<p>enum InfraredPort {</p>
<p>IF_PORT_GP1;</p>
<p>IF_PORT_GP2;</p>
<p>IF_PORT_GP4;</p>
<p>IF_PORT_GP5;</p>
<p>};</p>
<p>infraredPort：光电传感器连接至机械臂的接口，请选择对应的接口</p>
<p>value：传感器数值</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark116
> .anchor}![](media/image298.png){width="0.4066655730533683in"
> height="0.12999890638670167in"} **使能颜色传感器**
>
> 表 18.13 使能颜色传感器

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetColorSensor(bool enable,ColorPort colorPort, uint8_t
version)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>使能颜色传感器</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>ColorPort 定义：</p>
<p>enum ColorPort {</p>
<p>IF_PORT_GP1;</p>
<p>IF_PORT_GP2;</p>
<p>IF_PORT_GP4;</p>
<p>IF_PORT_GP5;</p>
<p>};</p>
<p>enable：使能标志。0：未使能。 1：使能</p>
<p>colorPort：颜色传感器连接至机械臂的接口，请选择对应的接口</p>
<p>version：颜色传感器版本号。0：版本为 V1.0 。1：版本为 V2.0。</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark117
> .anchor}![](media/image299.png){width="0.40999890638670167in"
> height="0.13in"} **获取颜色传感器读数**
>
> 表 18.14 获取颜色传感器读数

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetColorSensor( uint8_t *r, uint8_t *g, uint8_t *b)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取颜色传感器读数</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>r：红色。取值范围： 0-255</p>
<p>g：绿色。取值范围：0-255</p>
<p>b：蓝色。取值范围：0-255</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark118
> .anchor}![](media/image302.png){width="0.21666557305336834in"
> height="0.14333333333333334in"} **CAL 功能**
>
> []{#_bookmark119
> .anchor}![](media/image303.png){width="0.2949989063867017in"
> height="0.12999890638670167in"} **设置角度传感器静态偏差**
>
> 由于角度传感器焊接、机器状态等原因，大小臂上的角度传感器可能存在一个静态偏差。
> 我们可以通过各种手段（如调平、与标准源比较），得到此静态偏差，并通过此
> API 写入到 设备中。
>
> 大/小臂角度=大/小臂角度传感器静态偏差值+大/小臂角度传感器读数\*传感器线性化参
> 数
>
> 基座角度=基座编码器静态偏差值+基座编码器读数
>
> 表 19.1 设置角度传感器静态偏差

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetAngleSensorStaticError(float rearArmAngleError, float</p>
<p>frontArmAngleError)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置大小臂角度传感器静态偏差</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>rearArmAngleError：大臂角度传感器静态偏差</p>
<p>frontArmAngleError：小臂角度传感器静态偏差</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark120
> .anchor}![](media/image304.png){width="0.313332239720035in"
> height="0.12999890638670167in"} **读取角度传感器静态偏差**
>
> 表 19.2 读取角度传感器静态偏差

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetAngleSensorStaticError(float *rearArmAngleError, float</p>
<p>*frontArmAngleError)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>读取大/小臂角度传感器静态偏差</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>rearArmAngleError：大臂角度传感器静态偏差</p>
<p>frontArmAngleError：小臂角度传感器静态偏差</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark121 .anchor}![](media/image305.png){width="0.315in"
> height="0.13in"} **设置角度传感器线性化参数**
>
> 表 19.3 设置角度传感器线性化参数

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetAngleSensorCoef (float rearArmAngleCoef, float
frontArmAngleCoef)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置大/小臂角度传感器线性化参数</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>rearArmAngleCoef：大臂角度传感器线性化参数</p>
<p>frontArmAngleCoef：小臂角度传感器线性化参数</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark122
> .anchor}![](media/image307.png){width="0.318332239720035in"
> height="0.12999890638670167in"} **读取角度传感器线性化参数**
>
> 表 19.4 读取角度传感器线性化参数

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetAngleSensorCoef(float *rearArmAngleCoef, float
*frontArmAngleCoef)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>读取大/小臂角度传感器线性化参数</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>rearArmAngleCoef：大臂角度传感器线性化参数</p>
<p>frontArmAngleCoef：小臂角度传感器线性化参数</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark123
> .anchor}![](media/image308.png){width="0.31666557305336834in"
> height="0.13in"} **设置基座编码器静态偏差**
>
> 表 19.5 设置基座编码器静态偏差

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetBaseDecoderStaticError (float baseDecoderError)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置基座编码器静态偏差</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>baseDecoderError：基座编码器静态偏差</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark124
> .anchor}![](media/image309.png){width="0.31666557305336834in"
> height="0.12999890638670167in"} **读取基座编码器静态偏差**
>
> 表 19.6 读取基座编码器静态偏差

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetBaseDecoderStaticError(float *baseDecoderError)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>读取基座编码器静态偏差</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>baseDecoderError：基座编码器静态偏差</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark125
> .anchor}![](media/image311.png){width="0.23499890638670165in"
> height="0.14333333333333334in"} **WIFI 功能**
>
> DobotStudio可通过WIFI控制Dobot Maigican 。Dobot
> Maigican连接WIFI模块后，需
> 设置无线网络相关参数（IP地址、子网掩码、默认网关，使能WIFI等），使Dobot
> Magician 接入无线局域网。接入成功后Dobot
> Magician无需通过USB即可连接DobotStudio。
>
> []{#_bookmark126
> .anchor}![](media/image312.png){width="0.303332239720035in"
> height="0.12999890638670167in"} **使能** **WIFI**
>
> 表 20.1 设置 WIFI 配置模式接口说明

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetWIFIConfigMode(bool enable)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>使能 WIFI</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>Enable：0 ，未使能。 1，使能</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark127
> .anchor}![](media/image313.png){width="0.32166666666666666in"
> height="0.12999890638670167in"} **获取** **WIFI 状态**
>
> 表 20.2 获取 WIFI 状态

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetWIFIConfigMode(bool *isEnabled)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取当前 WIFI 状态</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>isEnabled：0 ，未使能。 1，使能</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark128
> .anchor}![](media/image314.png){width="0.323332239720035in"
> height="0.12999890638670167in"} **设置** **SSID**
>
> SSID（Service Set Identifier）：无线网络名称。
>
> 表 20.3 设置网络 SSID 接口说明

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetWIFISSID(const char *ssid)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置网络 SSID</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>ssid：字符串指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark129
> .anchor}![](media/image315.png){width="0.32666666666666666in"
> height="0.13in"} **获取设置的** **SSID**
>
> 表 20.4 获取当前设置 SSID

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetWIFISSID(char *ssid,uint32_tmaxLen)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取当前设置的 SSID</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>ssid：字符串指针</p>
<p>maxLen：字符串最大长度，以避免溢出</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark130 .anchor}![](media/image316.png){width="0.325in"
> height="0.12999890638670167in"} **设置** **WIFI 密码**
>
> 表 20.5 设置网络密码接口说明

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetWIFIPassword(const char *password)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置 WIFI 密码</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>password：字符串指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark131 .anchor}![](media/image317.png){width="0.325in"
> height="0.12999890638670167in"} **获取** **WIFI 密码**
>
> 表 20.6 获取当前设置网络密码

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetWIFIPassword(char *password,uint32_tmaxLen)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取 WIFI 密码</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>password：字符串指针</p>
<p>maxLen：字符串最大长度，以避免溢出</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark132
> .anchor}![](media/image318.png){width="0.323332239720035in"
> height="0.13in"} **设置** **IP 地址**
>
> 表 20.7 设置 IP 地址

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetWIFIIPAddress(WIFIIPAddress *wifiIPAddress)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置 IP 地址</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>typedef struct tagWIFIIPAddress {</p>
</blockquote>
<table>
<colgroup>
<col style="width: 28%" />
<col style="width: 71%" />
</colgroup>
<tbody>
<tr class="odd">
<td>uint8_t dhcp;</td>
<td><blockquote>
<p>//是否开启 DHCP 。0：关闭。1：开启</p>
</blockquote></td>
</tr>
<tr class="even">
<td>uint8_taddr[4];</td>
<td>//IP 地址分成四段， 每段取值范围：0~255</td>
</tr>
</tbody>
</table>
<blockquote>
<p>}WIFIIPAddress;</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td></td>
<td><blockquote>
<p>wifiIPAddr：WIFIIPAddress 指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark133
> .anchor}![](media/image319.png){width="0.323332239720035in"
> height="0.12999890638670167in"} **获取设置的** **IP 地址**
>
> 表 20.8 获取当前设置 IP 地址接口说明

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetWIFIIPAddress(WIFIIPAddress *wifiIPAddress)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取设置的 IP 地址</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>typedef struct tagWIFIIPAddress {</p>
<p>uint8_t dhcp; //是否开启 DHCP 。0：关闭。 1：开启</p>
<p>uint8_taddr[4]; //IP 地址分成四段， 每段取值范围：0~255</p>
<p>}WIFIIPAddress;</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark134
> .anchor}![](media/image320.png){width="0.32166666666666666in"
> height="0.12999890638670167in"} **设置子网掩码**
>
> 表 20.9 设置子网掩码

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetWIFINetmask(WIFINetmask *wifiNetmask)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置子网掩码</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>typedef struct tagWIFINetmask {</p>
<p>uint8_taddr[4]; //IP 地址分成四段， 每段取值范围：0~255</p>
<p>}WIFINetmask;</p>
<p>wifiNetmask：WIFINetmask 指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回，导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark135
> .anchor}![](media/image321.png){width="0.41333333333333333in"
> height="0.13in"} **获取设置的子网掩码**
>
> 表 20.10 获取当前设置子网掩码接口说明

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetWIFINetmask(WIFINetmask *wifiNetmask)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取设置的子网掩码</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>typedef struct tagWIFINetmask {</p>
<p>uint8_t addr[4]; //IP 地址分成四段， 每段取值范围：0~255</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td></td>
<td><blockquote>
<p>}WIFINetmask;</p>
<p>wifiNetmask：WIFINetmask 指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark136 .anchor}![](media/image322.png){width="0.395in"
> height="0.12999890638670167in"} **设置网关**
>
> 表 20.11 设置网关

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetWIFIGateway(WIFIGateway *wifiGateway)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置网关</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>typedef struct tagWIFIGateway {</p>
<p>uint8_t addr[4]; //IP地址分成四段，每段取值范：0~255</p>
<p>}WIFIGateway;</p>
<p>wifiGateway：WIFIGateway指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark137
> .anchor}![](media/image323.png){width="0.41333333333333333in"
> height="0.13in"} **获取设置的网关**
>
> 表 20.12 获取当前设置子网掩码接口说明

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetWIFIGateway(WIFIGateway *wifiGateway)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取设置的网关</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>typedef struct tagWIFIGateway {</p>
<p>uint8_t addr[4]; //IP地址分成四段，每段取值范：0~255</p>
<p>}WIFIGateway;</p>
<p>wifiGateway：WIFIGateway 指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark138
> .anchor}![](media/image324.png){width="0.41499890638670167in"
> height="0.13in"} **设置** **DNS**
>
> 表 20.13 设置 DNS

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int SetWIFIDNS(WIFIDNS *wifiDNS)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>设置 DNS</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>typedef struct tagWIFIDNS {</p>
<p>uint8_t addr[4]; //IP地址分成四段，每段取值范：0~255</p>
</blockquote></td>
</tr>
</tbody>
</table>

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td></td>
<td><blockquote>
<p>}WIFIDNS;</p>
<p>wifiDNS：WIFIDNS 指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark139
> .anchor}![](media/image325.png){width="0.41833333333333333in"
> height="0.12999890638670167in"} **获取设置的** **DNS**
>
> 表 20.14 获取设置的 DNS

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetWIFIDNS(WIFIDNS *wifiDNS)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取设置的 DNS</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>typedef struct tagWIFIDNS {</p>
<p>uint8_taddr[4]; //IP 分成四段， 每段取值范围 0~255</p>
<p>}WIFIDNS;</p>
<p>wifiDNS：DNS 结构体指针</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark140
> .anchor}![](media/image326.png){width="0.4166666666666667in"
> height="0.13in"} **获取当前** **WIFI 模块的连接状态**
>
> 表 20.15 获取当前 WIFI 模块连接状态

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>int GetWIFIConnectStatus(bool *isConnected)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>获取当前 WIFI 模块连接状态</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>isConnected：0 ，未连接。1，连接</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>DobotCommunicate_NoError：指令正常返回</p>
<p>DobotCommunicate_Timeout：指令无返回， 导致超时</p>
</blockquote></td>
</tr>
</tbody>
</table>

> []{#_bookmark141
> .anchor}![](media/image328.png){width="0.23499890638670165in"
> height="0.14333333333333334in"} **其他功能**
>
> ![](media/image329.png){width="0.303332239720035in"
> height="0.12666557305336834in"} **事件循环功能**
>
> 在某些语言中， 当调用 API 接口后， 如果没有事件循环，
> 应用程序将直接退出，导致指 令没有下发至 Dobot
> 控制器。为避免这种情况发生，我们提供了事件循环接口， 在应用程序
> 退出前调用（目前已知需要做此处理的语言有 Python）。
>
> 表 21.1 事件循环功能接口说明

<table>
<colgroup>
<col style="width: 9%" />
<col style="width: 90%" />
</colgroup>
<tbody>
<tr class="odd">
<td><blockquote>
<p>原型</p>
</blockquote></td>
<td><blockquote>
<p>void DobotExec(void)</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>功能</p>
</blockquote></td>
<td><blockquote>
<p>事件循环功能</p>
</blockquote></td>
</tr>
<tr class="odd">
<td><blockquote>
<p>参数</p>
</blockquote></td>
<td><blockquote>
<p>无</p>
</blockquote></td>
</tr>
<tr class="even">
<td><blockquote>
<p>返回</p>
</blockquote></td>
<td><blockquote>
<p>无</p>
</blockquote></td>
</tr>
</tbody>
</table>
