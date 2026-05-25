# Wails + Vue 3 桌面应用开发示例

## 项目简介

本项目是一个基于 **Wails v2** + **Vue 3** 的桌面应用示例，演示了 Go 后端与前端之间的多种交互模式。涵盖了字符串传参、对象传递、列表数据管理、窗口控制以及错误处理等常见场景。

## 技术栈

| 层级 | 技术 | 说明 |
|------|------|------|
| 桌面框架 | Wails v2.12.0 | 使用 Go 驱动原生 WebView，替代 Electron |
| 后端语言 | Go 1.23+ | 负责业务逻辑、状态管理和系统调用 |
| 前端框架 | Vue 3 (Composition API) | 使用 `<script setup>` 语法 |
| 构建工具 | Vite 3 | 开发服务器与生产打包 |

## 项目结构

```
wails/
├── app.go                    # Go 后端核心代码（数据结构与方法定义）
├── main.go                   # 应用入口，Wails 配置与菜单
├── go.mod / go.sum           # Go 模块依赖
├── wails.json                # Wails 项目配置
├── build/                    # 构建输出目录
│   ├── appicon.png           # 应用图标
│   └── darwin/               # macOS 平台构建产物
├── frontend/                 # 前端项目根目录
│   ├── index.html            # HTML 入口
│   ├── package.json          # 前端依赖（Vue 3 + Vite 3）
│   ├── vite.config.js        # Vite 构建配置
│   ├── src/
│   │   ├── main.js           # Vue 应用入口
│   │   ├── App.vue           # 根组件
│   │   ├── style.css         # 全局样式
│   │   ├── components/
│   │   │   └── HelloWorld.vue # 主交互组件（核心演示页面）
│   │   └── assets/
│   │       ├── images/
│   │       │   └── logo-universal.png
│   │       └── fonts/
│   │           └── nunito-v16-latin-regular.woff2
│   └── wailsjs/              # Wails 自动生成的绑定代码
│       └── go/main/App.js    # Go 方法的前端调用代理
└── README.md                 # 英文原始说明
```

## 快速开始

### 环境要求

- **Go** 1.23 或更高版本
- **Node.js** 16 或更高版本
- **Wails CLI** — 安装方式：

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

- **macOS**: Xcode Command Line Tools
- **Windows**: WebView2 Runtime（Windows 10/11 通常已内置）
- **Linux**: `gtk3`、`webkit2gtk` 等系统库

### 开发模式

在项目根目录执行：

```bash
wails dev
```

该命令会同时启动：
- Vite 开发服务器（前端热更新）
- Go 后端运行时（通过 WebSocket 与前端通信）

开发服务器同时暴露一个浏览器地址（默认 `http://localhost:34115`），可以在浏览器中打开并用 DevTools 调试前端，同时仍可调用 Go 方法。

### 构建生产版本

```bash
wails build
```

构建产物输出到 `build/` 目录，生成对应平台的原生可执行文件。

### 基本命令速查

| 命令 | 说明 |
|------|------|
| `wails dev` | 启动开发模式（热更新） |
| `wails build` | 构建生产包 |
| `wails build -platform windows` | 交叉编译 Windows 版本 |
| `wails doctor` | 检查开发环境配置 |

## 后端 API 详解

所有后端方法定义在 `app.go`，通过 `main.go` 中的 `options.App.Bind` 绑定后，Wails 会自动为这些方法生成前端调用代理（位于 `frontend/wailsjs/go/main/App.js`）。

### 启动流程

```
main.go:
  NewApp()             → 创建 App 实例
  wails.Run(...)       → 启动 Wails 运行时
    └─ OnStartup        → 调用 app.startup(ctx)，保存上下文
    └─ Bind: [app]      → 将 App 的所有公开方法绑定为前端可调用 API
```

### 核心数据结构

#### ProfileInput — 传给 Go 的对象参数

```go
type ProfileInput struct {
    Name  string `json:"name"`  // 姓名
    Role  string `json:"role"`  // 岗位
    Years int    `json:"years"` // 工作年限
}
```

#### ProfileSummary — Go 返回给前端的对象

```go
type ProfileSummary struct {
    Title   string   `json:"title"`   // 摘要标题（如 "李雷 / 调度员"）
    Message string   `json:"message"` // 描述信息
    Score   int      `json:"score"`   // 评分 (0-100)
    Tags    []string `json:"tags"`    // 标签列表
}
```

#### TaskInput — 新增任务入参

```go
type TaskInput struct {
    Title    string `json:"title"`    // 任务标题
    Priority string `json:"priority"` // 优先级（普通/重要/紧急）
}
```

#### Task — 任务实体

```go
type Task struct {
    ID        int    `json:"id"`        // 自增 ID
    Title     string `json:"title"`     // 任务标题
    Priority  string `json:"priority"`  // 优先级
    Done      bool   `json:"done"`      // 是否完成
    CreatedAt string `json:"createdAt"` // 创建时间
}
```

#### InteractionState — 后端状态快照

```go
type InteractionState struct {
    CallCount   int    `json:"callCount"`   // 累计调用次数
    LastMessage string `json:"lastMessage"` // 最近一次操作描述
    Tasks       []Task `json:"tasks"`       // 任务列表
}
```

### API 方法一览

#### 1. Greet(name string) string

最简单的字符串传参示例。

- **入参**: 姓名字符串
- **返回**: 格式化的问候语
- **前端调用**:
  ```js
  const result = await Greet("李雷")
  // → "你好 李雷，Go 后端已经收到这个字符串参数。"
  ```

#### 2. BuildProfile(input ProfileInput) (ProfileSummary, error)

对象作为参数，Go 计算后返回结构化结果，同时演示**错误返回**。

- **入参**: `{ name, role, years }`
- **返回**: `{ title, message, score, tags }`
- **评分规则**: 基础分 60 + 工作年限 × 8，上限 100
- **校验**: `name` 为空时返回 error
- **前端调用**:
  ```js
  const result = await BuildProfile({ name: "李雷", role: "调度员", years: 3 })
  // → { title: "李雷 / 调度员", message: "后端根据 3 年经验计算出示例评分 84。", score: 84, tags: [...] }
  ```

#### 3. AddTask(input TaskInput) (Task, error)

新增任务，写入 Go 侧内存状态，演示**带 ID 和时间戳的写操作**。

- **入参**: `{ title, priority }`
- **返回**: 完整的 Task 对象（含自动生成的 ID 和创建时间）
- **校验**: `title` 为空时返回 error
- **注意**: 新增的任务插入到列表头部

#### 4. GetState() InteractionState

获取后端当前状态的完整快照，用于前端刷新展示。

- **返回**: `{ callCount, lastMessage, tasks }`
- **线程安全**: 使用 `sync.Mutex` 保护并发访问
- **使用场景**: 每次调用其他方法后自动拉取，也支持手动刷新

#### 5. 窗口控制方法组

| 方法 | 说明 | 返回 |
|------|------|------|
| `ShowWindow()` | 恢复窗口显示 | `WindowCommandResult` |
| `HideToTray()` | 隐藏窗口（模拟后台运行） | `WindowCommandResult` |
| `CenterWindow()` | 窗口居中到屏幕 | `WindowCommandResult` |
| `SetWindowTitle(title)` | 修改窗口标题 | `WindowCommandResult` |
| `SetWindowSize({width, height})` | 修改窗口尺寸（最小 640×480） | `WindowCommandResult` |
| `QuitApp()` | 退出应用 | `error` |

窗口方法依赖 `startup` 阶段保存的 `context.Context`，调用前会校验上下文就绪状态。尺寸设置有最低限制：宽度 ≥ 640，高度 ≥ 480。

#### 6. GetTrayStatus() TrayStatus

返回当前托盘式后台运行模式的说明信息，非真正的系统托盘 API。

```js
const status = await GetTrayStatus()
// → { mode: "托盘式后台运行", supported: false, notes: [...], menuItems: [...] }
```

### 并发安全

`App` 结构体使用 `sync.Mutex` 保护以下共享字段：

- `callCount` — 调用计数器
- `lastMessage` — 最后操作描述
- `tasks` — 任务切片（返回时使用 `copy` 创建快照，避免并发读写冲突）
- `nextTaskID` — 自增 ID 计数器

## 前端架构

### 组件结构

```
App.vue
  └── HelloWorld.vue    （全部交互逻辑集中于此组件）
```

### HelloWorld.vue 核心设计

组件围绕"一个后端调用收敛函数"构建，所有 Go 调用统一走 `runBackend()`：

```js
async function runBackend(key, action) {
    if (loading[key]) return            // 防止重复请求
    loading[key] = true
    errorText.value = ''
    try {
        await action()
    } catch (error) {
        errorText.value = error?.message || String(error)
    } finally {
        loading[key] = false
    }
}
```

**设计要点**:
- 每个操作类型有独立的 `loading` 状态（`loading.greet`、`loading.profile` 等）
- 全局错误统一展示在 `.error-text` 区域
- 每次操作后自动调用 `loadState()` 刷新状态快照
- 使用 Vue 3 的 `reactive` 管理表单数据，`ref` 管理展示数据

### 页面功能分区

| 区域 | 对应方法 | 演示内容 |
|------|----------|----------|
| 字符串参数 | `Greet(name)` | 简单字符串入参/出参 |
| 对象传递 | `BuildProfile(input)` | 结构化对象传参、计算、标签生成 |
| 窗口控制 | `ShowWindow` / `HideToTray` / `CenterWindow` / `SetWindowTitle` / `SetWindowSize` | 系统窗口 API 调用 |
| 列表和状态 | `AddTask(input)` / `GetState()` | 列表 CRUD、状态同步 |
| 状态快照 | `GetState()` 结果展示 | JSON 格式查看后端完整状态 |

### Go 方法的前端导入方式

```js
import {
  Greet,
  BuildProfile,
  AddTask,
  GetState,
  ShowWindow,
  HideToTray,
  CenterWindow,
  SetWindowTitle,
  SetWindowSize,
  GetTrayStatus,
} from '../../wailsjs/go/main/App'
```

这些是由 Wails 构建工具链自动生成的代理函数，调用时会通过 WebSocket 将参数序列化到 Go 端执行，再将结果反序列化回前端。

## 运行原理

### Wails 通信模型

```
┌─────────────────────────────────────────────┐
│  Wails 桌面应用                              │
│                                             │
│  ┌──────────┐      WebSocket      ┌───────┐ │
│  │ Go 后端  │ ◄──────────────────► │ WebView│ │
│  │          │    JSON 序列化/反序列化  │ (Vue) │ │
│  │ app.go   │                      │       │ │
│  └──────────┘                      └───────┘ │
│       │                                   │  │
│       │  Wails Runtime API                │  │
│       ├─ WindowShow / WindowHide          │  │
│       ├─ WindowSetTitle / WindowSetSize   │  │
│       ├─ WindowCenter                     │  │
│       └─ Quit                             │  │
└─────────────────────────────────────────────┘
```

1. 前端调用代理函数 → Wails 将参数序列化为 JSON
2. JSON 通过 WebSocket 发送到 Go 后端
3. Go 方法执行（可能调用 Wails Runtime API 控制系统窗口）
4. 返回值序列化为 JSON → 通过 WebSocket 返回前端
5. 前端接收结果并更新 UI

### 菜单系统

`main.go` 中构建了跨平台应用菜单：

- **macOS**: `App 菜单` + `编辑菜单` + `窗口菜单` + `窗口控制`（自定义）
- **Windows**: `编辑菜单` + `窗口菜单` + `窗口控制`（自定义）

自定义"窗口控制"子菜单包含：显示窗口、隐藏到后台、窗口居中、退出应用。这些菜单回调直接复用 `App` 的方法，与前端调用共享同一套后端逻辑。

### 托盘式后台运行

当前示例采用"窗口隐藏 + 菜单恢复"方案模拟后台运行：

1. `wails.json` 中未配置系统托盘（Wails v2.12 无稳定跨平台托盘 API）
2. `HideWindowOnClose: true` — 点关闭按钮时隐藏窗口而非退出
3. 通过菜单"显示窗口"或前端按钮恢复
4. `GetTrayStatus()` 向前端说明当前方案的局限性

## 开发注意事项

### 数据传递限制

- 前后端通过 JSON 序列化通信，Go 方法的参数和返回值必须可被 `encoding/json` 处理
- 目前不支持流式传输、文件上传等高级模式
- 复杂对象建议定义专用结构体，确保 JSON tag 与前端字段名一致

### 上下文就绪检查

所有调用 Wails Runtime API 的方法都应先检查 `a.ctx` 是否已由 `startup` 初始化。示例中通过 `appContext()` 辅助方法统一处理：

```go
func (a *App) appContext() (context.Context, error) {
    if a.ctx == nil {
        return nil, errors.New("Wails 上下文尚未就绪")
    }
    return a.ctx, nil
}
```

### 跨平台构建

如果需要在 macOS 上交叉编译 Windows 版本：

```bash
wails build -platform windows
```

Windows 上交叉编译 macOS 版本需要 Darwin 工具链支持，通常不推荐。

### Binding 规则

`main.go` 中 `Bind` 切片内的所有公开（首字母大写）方法都会被暴露给前端。私有方法不会被暴露。本项目中只有 `App` 实例被绑定，其所有公开方法均在前端可见。

## 常见问题

### Q: 前端报 "Go 方法未定义"？

确认：
1. `wails dev` 或 `wails build` 已重新执行（会更新 `wailsjs/` 下的绑定代码）
2. 方法名首字母大写（Go 的导出规则）
3. 方法已通过 `Bind` 注册

### Q: 调用窗口方法时报 "上下文尚未就绪"？

这是正常的启动时序保护。`startup` 在 Wails 运行时就绪后被调用，此后 `a.ctx` 才有效。确保不在 `init()` 或包级变量初始化阶段调用需要上下文的方法。

### Q: 能不能在浏览器中调试？

可以。`wails dev` 启动后会在 `http://localhost:34115` 提供浏览器可访问的版本。在浏览器 DevTools 中打开该地址即可调试前端，Go 方法调用依然有效（通过 WebSocket 回退到 localhost）。

## 相关资源

- [Wails 官方文档](https://wails.io/docs/introduction)
- [Wails v2 项目配置参考](https://wails.io/docs/reference/project-config)
- [Vue 3 文档](https://cn.vuejs.org/guide/introduction.html)
- [Vite 文档](https://cn.vitejs.dev/)
