/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/8/7 21:50
 * @version 1.0
 */
package taskx

import (
	"fmt"
	"time"
)

/*
*
定时任务-异常不自动恢复
任务执行开始延时时间，单位秒
任务执行完毕延时时间，单位秒
具体任务
*/
func StartTaskNoRecover(delaySecs, interSecs int, cmd func()) {
	go func() {
		if delaySecs > 0 {
			time.Sleep(time.Duration(delaySecs) * time.Second)
		}
		for {
			// 执行定时任务
			cmd()
			time.Sleep(time.Duration(interSecs) * time.Second)
		}
	}()
}

/*
*
定时任务-异常自动恢复
任务执行开始延时时间，单位秒
任务执行完毕延时时间，单位秒
具体任务
*/
func StartTask(delaySecs, interSecs int, cmd func()) {
	go func() {
		if delaySecs > 0 {
			time.Sleep(time.Duration(delaySecs) * time.Second)
		}
		for {
			// 异常自动恢复
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered in defer", r)
				}
			}()
			// 执行定时任务
			cmd()
			time.Sleep(time.Duration(interSecs) * time.Second)
		}
	}()
}
