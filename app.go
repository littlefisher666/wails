package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App 保存桌面应用运行期状态，示例数据放在 Go 侧用于演示跨调用共享。
type App struct {
	ctx         context.Context
	mu          sync.Mutex
	callCount   int
	lastMessage string
	nextTaskID  int
	tasks       []Task
}

// ProfileInput 是前端传给 Go 的对象参数示例。
type ProfileInput struct {
	Name  string `json:"name"`
	Role  string `json:"role"`
	Years int    `json:"years"`
}

// ProfileSummary 是 Go 处理后返回给前端的对象结果示例。
type ProfileSummary struct {
	Title   string   `json:"title"`
	Message string   `json:"message"`
	Score   int      `json:"score"`
	Tags    []string `json:"tags"`
}

// TaskInput 是新增任务时的前端输入。
type TaskInput struct {
	Title    string `json:"title"`
	Priority string `json:"priority"`
}

// Task 是列表数据往返传递的示例。
type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Priority  string `json:"priority"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"createdAt"`
}

// InteractionState 让前端一次拿到后端当前状态快照。
type InteractionState struct {
	CallCount   int    `json:"callCount"`
	LastMessage string `json:"lastMessage"`
	Tasks       []Task `json:"tasks"`
}

// WindowSizeInput 是前端设置窗口尺寸时传给 Go 的对象。
type WindowSizeInput struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// WindowCommandResult 统一返回窗口命令结果，方便前端展示执行反馈。
type WindowCommandResult struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

// TrayStatus 描述当前 Wails v2 示例采用的托盘式后台运行方式。
type TrayStatus struct {
	Mode      string   `json:"mode"`
	Supported bool     `json:"supported"`
	Notes     []string `json:"notes"`
	MenuItems []string `json:"menuItems"`
}

// NewApp 创建应用实例，并准备一条初始任务方便前端立即展示列表结构。
func NewApp() *App {
	return &App{
		nextTaskID: 2,
		tasks: []Task{
			{
				ID:        1,
				Title:     "查看 Go 方法如何返回列表数据",
				Priority:  "普通",
				CreatedAt: time.Now().Format(time.DateTime),
			},
		},
	}
}

// startup 保存 Wails 上下文，后续需要调用 runtime API 时可以复用。
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet 接收前端传来的字符串，并返回一个字符串结果。
func (a *App) Greet(name string) string {
	cleanName := strings.TrimSpace(name)
	if cleanName == "" {
		cleanName = "Wails"
	}

	message := fmt.Sprintf("你好 %s，Go 后端已经收到这个字符串参数。", cleanName)
	a.recordCall(message)
	return message
}

// BuildProfile 接收前端对象，返回经过 Go 侧计算后的对象。
func (a *App) BuildProfile(input ProfileInput) (ProfileSummary, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return ProfileSummary{}, errors.New("姓名不能为空")
	}

	role := strings.TrimSpace(input.Role)
	if role == "" {
		role = "未设置岗位"
	}

	years := input.Years
	if years < 0 {
		years = 0
	}

	score := 60 + years*8
	if score > 100 {
		score = 100
	}

	summary := ProfileSummary{
		Title:   fmt.Sprintf("%s / %s", name, role),
		Message: fmt.Sprintf("后端根据 %d 年经验计算出示例评分 %d。", years, score),
		Score:   score,
		Tags:    buildProfileTags(role, years),
	}

	a.recordCall("已生成资料摘要：" + summary.Title)
	return summary, nil
}

// AddTask 接收前端对象并写入 Go 侧状态，返回带 ID 和时间的任务对象。
func (a *App) AddTask(input TaskInput) (Task, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return Task{}, errors.New("任务内容不能为空")
	}

	priority := strings.TrimSpace(input.Priority)
	if priority == "" {
		priority = "普通"
	}

	a.mu.Lock()
	task := Task{
		ID:        a.nextTaskID,
		Title:     title,
		Priority:  priority,
		CreatedAt: time.Now().Format(time.DateTime),
	}
	a.nextTaskID++
	a.tasks = append([]Task{task}, a.tasks...)
	a.callCount++
	a.lastMessage = "已新增任务：" + task.Title
	a.mu.Unlock()

	return task, nil
}

// GetState 返回当前后端状态快照，前端可用它刷新列表和调试数据。
func (a *App) GetState() InteractionState {
	a.mu.Lock()
	defer a.mu.Unlock()

	return InteractionState{
		CallCount:   a.callCount,
		LastMessage: a.lastMessage,
		Tasks:       a.taskSnapshotLocked(),
	}
}

// ShowWindow 从后台恢复窗口显示，主要配合菜单里的“显示窗口”使用。
func (a *App) ShowWindow() (WindowCommandResult, error) {
	return a.runWindowCommand("显示窗口", "窗口已恢复显示。", wruntime.WindowShow)
}

// HideToTray 隐藏窗口，用来演示托盘式后台运行入口。
func (a *App) HideToTray() (WindowCommandResult, error) {
	return a.runWindowCommand("隐藏到后台", "窗口已隐藏，可通过应用菜单恢复。", wruntime.WindowHide)
}

// CenterWindow 将窗口移动到当前屏幕中央。
func (a *App) CenterWindow() (WindowCommandResult, error) {
	return a.runWindowCommand("居中窗口", "窗口已移动到屏幕中央。", wruntime.WindowCenter)
}

// SetWindowTitle 修改原生窗口标题，演示前端字符串传到 Go 后再调用 runtime。
func (a *App) SetWindowTitle(title string) (WindowCommandResult, error) {
	cleanTitle := strings.TrimSpace(title)
	if cleanTitle == "" {
		cleanTitle = "helloworld"
	}

	return a.runWindowCommand("修改标题", "窗口标题已改为："+cleanTitle, func(ctx context.Context) {
		wruntime.WindowSetTitle(ctx, cleanTitle)
	})
}

// SetWindowSize 修改原生窗口尺寸，演示前端对象参数驱动窗口 runtime。
func (a *App) SetWindowSize(input WindowSizeInput) (WindowCommandResult, error) {
	width := input.Width
	if width < 640 {
		width = 640
	}

	height := input.Height
	if height < 480 {
		height = 480
	}

	return a.runWindowCommand("修改尺寸", fmt.Sprintf("窗口尺寸已设置为 %d x %d。", width, height), func(ctx context.Context) {
		wruntime.WindowSetSize(ctx, width, height)
	})
}

// QuitApp 退出应用，菜单和前端都可以复用这条后端能力。
func (a *App) QuitApp() error {
	ctx, err := a.appContext()
	if err != nil {
		return err
	}

	wruntime.Quit(ctx)
	return nil
}

// GetTrayStatus 返回托盘式示例说明，避免把 Wails v2 的能力说成完整系统托盘。
func (a *App) GetTrayStatus() TrayStatus {
	return TrayStatus{
		Mode:      "托盘式后台运行",
		Supported: false,
		Notes: []string{
			"Wails v2.12 没有 options.App 级别的稳定跨平台系统托盘配置。",
			"当前示例使用关闭隐藏、窗口隐藏和应用菜单恢复，先覆盖后台运行的最短链路。",
			"如需真正系统托盘图标和右键菜单，建议升级到 Wails v3 托盘 API 或接专门托盘库。",
		},
		MenuItems: []string{"显示窗口", "隐藏到后台", "窗口居中", "退出应用"},
	}
}

func (a *App) recordCall(message string) {
	a.mu.Lock()
	a.callCount++
	a.lastMessage = message
	a.mu.Unlock()
}

func (a *App) taskSnapshotLocked() []Task {
	// 返回快照，避免前端序列化期间读到正在被新增任务修改的切片。
	tasks := make([]Task, len(a.tasks))
	copy(tasks, a.tasks)
	return tasks
}

func (a *App) appContext() (context.Context, error) {
	if a.ctx == nil {
		return nil, errors.New("Wails 上下文尚未就绪")
	}
	return a.ctx, nil
}

func (a *App) runWindowCommand(action string, message string, command func(context.Context)) (WindowCommandResult, error) {
	ctx, err := a.appContext()
	if err != nil {
		return WindowCommandResult{}, err
	}

	command(ctx)
	a.recordCall(message)
	return WindowCommandResult{
		Action:  action,
		Message: message,
	}, nil
}

func buildProfileTags(role string, years int) []string {
	tags := []string{"对象参数", "Go 计算"}
	if years >= 3 {
		tags = append(tags, "经验充足")
	}
	if strings.Contains(role, "调度") {
		tags = append(tags, "调度岗位")
	}
	return tags
}
