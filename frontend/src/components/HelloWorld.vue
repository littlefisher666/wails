<script setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {
  AddTask,
  BuildProfile,
  CenterWindow,
  GetState,
  GetTrayStatus,
  Greet,
  HideToTray,
  SetWindowSize,
  SetWindowTitle,
  ShowWindow,
} from '../../wailsjs/go/main/App'

const profileForm = reactive({
  name: '李雷',
  role: '调度员',
  years: 3,
})

const taskForm = reactive({
  title: '巡检主运输皮带运行状态',
  priority: '重要',
})

const windowForm = reactive({
  title: 'Wails 控制台示例',
  width: 1024,
  height: 768,
})

const priorities = ['普通', '重要', '紧急']
const greeting = ref('等待调用 Go 后端')
const profileSummary = ref(null)
const trayStatus = ref(null)
const windowMessage = ref('等待窗口命令')
const errorText = ref('')

const state = reactive({
  callCount: 0,
  lastMessage: '',
  tasks: [],
})

const loading = reactive({
  greet: false,
  profile: false,
  task: false,
  state: false,
  window: false,
  tray: false,
})

const stateJson = computed(() => JSON.stringify({
  callCount: state.callCount,
  lastMessage: state.lastMessage,
  tasks: state.tasks,
}, null, 2))

// 统一收敛 loading 和错误，让每个示例只关注自己的数据形态。
async function runBackend(key, action) {
  if (loading[key]) {
    return
  }

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

async function runWindowCommand(command) {
  await runBackend('window', async () => {
    const result = await command()
    windowMessage.value = result.message
    await loadState()
  })
}

async function submitGreeting() {
  await runBackend('greet', async () => {
    greeting.value = await Greet(profileForm.name)
    await loadState()
  })
}

async function submitProfile() {
  await runBackend('profile', async () => {
    profileSummary.value = await BuildProfile({
      name: profileForm.name,
      role: profileForm.role,
      years: Number(profileForm.years) || 0,
    })
    await loadState()
  })
}

async function submitTask() {
  await runBackend('task', async () => {
    await AddTask({
      title: taskForm.title,
      priority: taskForm.priority,
    })
    taskForm.title = ''
    await loadState()
  })
}

async function loadState() {
  await runBackend('state', async () => {
    const nextState = await GetState()
    state.callCount = nextState.callCount
    state.lastMessage = nextState.lastMessage
    state.tasks = nextState.tasks || []
  })
}

async function loadTrayStatus() {
  await runBackend('tray', async () => {
    trayStatus.value = await GetTrayStatus()
  })
}

async function applyWindowTitle() {
  await runWindowCommand(() => SetWindowTitle(windowForm.title))
}

async function applyWindowSize() {
  await runWindowCommand(() => SetWindowSize({
    width: Number(windowForm.width) || 1024,
    height: Number(windowForm.height) || 768,
  }))
}

const windowActions = [
  {label: '显示窗口', run: () => runWindowCommand(ShowWindow)},
  {label: '隐藏到后台', run: () => runWindowCommand(HideToTray)},
  {label: '窗口居中', run: () => runWindowCommand(CenterWindow)},
]

// 首屏主动拉一次后端状态，方便验证 Go 侧状态是否已和 Vue 连通。
onMounted(() => {
  loadState()
  loadTrayStatus()
})
</script>

<template>
  <main class="demo-shell">
    <section class="hero">
      <div class="hero-copy">
        <p class="eyebrow">Wails + Vue</p>
        <h1>前后端交互示例</h1>
        <p class="intro">这个页面演示 Vue 通过 Wails 绑定调用 Go 方法，包含字符串参数、对象传递、列表状态和错误返回。</p>
      </div>
      <img class="logo" alt="Wails logo" src="../assets/images/logo-universal.png">
    </section>

    <section class="status-strip" aria-label="后端状态">
      <div>
        <span>调用次数</span>
        <strong>{{ state.callCount }}</strong>
      </div>
      <div>
        <span>最后消息</span>
        <strong>{{ state.lastMessage || '暂无调用' }}</strong>
      </div>
      <button class="ghost-button" :disabled="loading.state" @click="loadState">
        刷新状态
      </button>
    </section>

    <p v-if="errorText" class="error-text">{{ errorText }}</p>

    <section class="grid">
      <article class="panel">
        <p class="panel-kicker">字符串参数</p>
        <h2>调用 Greet(name)</h2>
        <label>
          姓名
          <input v-model="profileForm.name" class="field" type="text" autocomplete="off">
        </label>
        <button class="primary-button" :disabled="loading.greet" @click="submitGreeting">
          {{ loading.greet ? '调用中' : '发送字符串' }}
        </button>
        <output class="result">{{ greeting }}</output>
      </article>

      <article class="panel">
        <p class="panel-kicker">对象传递</p>
        <h2>调用 BuildProfile(input)</h2>
        <label>
          岗位
          <input v-model="profileForm.role" class="field" type="text" autocomplete="off">
        </label>
        <label>
          年限
          <input v-model.number="profileForm.years" class="field" type="number" min="0">
        </label>
        <button class="primary-button" :disabled="loading.profile" @click="submitProfile">
          {{ loading.profile ? '生成中' : '生成资料摘要' }}
        </button>
        <div v-if="profileSummary" class="summary">
          <strong>{{ profileSummary.title }}</strong>
          <span>{{ profileSummary.message }}</span>
          <span>评分：{{ profileSummary.score }}</span>
          <div class="tag-row">
            <span v-for="tag in profileSummary.tags" :key="tag" class="tag">{{ tag }}</span>
          </div>
        </div>
      </article>

      <article class="panel panel-wide">
        <p class="panel-kicker">窗口和后台</p>
        <h2>窗口控制 / 托盘式后台</h2>
        <div class="button-row">
          <button v-for="action in windowActions" :key="action.label" class="ghost-button" :disabled="loading.window" @click="action.run">
            {{ action.label }}
          </button>
        </div>
        <div class="task-form window-form">
          <label>
            窗口标题
            <input v-model="windowForm.title" class="field" type="text" autocomplete="off">
          </label>
          <label>
            宽度
            <input v-model.number="windowForm.width" class="field" type="number" min="640">
          </label>
          <label>
            高度
            <input v-model.number="windowForm.height" class="field" type="number" min="480">
          </label>
        </div>
        <div class="button-row">
          <button class="primary-button" :disabled="loading.window" @click="applyWindowTitle">
            修改标题
          </button>
          <button class="primary-button" :disabled="loading.window" @click="applyWindowSize">
            应用尺寸
          </button>
        </div>
        <output class="result">{{ windowMessage }}</output>
        <div v-if="trayStatus" class="summary">
          <strong>{{ trayStatus.mode }}</strong>
          <span>{{ trayStatus.notes.join(' ') }}</span>
          <div class="tag-row">
            <span v-for="item in trayStatus.menuItems" :key="item" class="tag">{{ item }}</span>
          </div>
        </div>
      </article>

      <article class="panel panel-wide">
        <p class="panel-kicker">列表和状态</p>
        <h2>调用 AddTask(input) / GetState()</h2>
        <div class="task-form">
          <label>
            任务内容
            <input v-model="taskForm.title" class="field" type="text" autocomplete="off">
          </label>
          <label>
            优先级
            <select v-model="taskForm.priority" class="field">
              <option v-for="priority in priorities" :key="priority" :value="priority">
                {{ priority }}
              </option>
            </select>
          </label>
          <button class="primary-button" :disabled="loading.task" @click="submitTask">
            {{ loading.task ? '保存中' : '新增任务' }}
          </button>
        </div>
        <ul class="task-list">
          <li v-for="task in state.tasks" :key="task.id">
            <span class="task-title">{{ task.title }}</span>
            <span>{{ task.priority }}</span>
            <time>{{ task.createdAt }}</time>
          </li>
        </ul>
      </article>

      <article class="panel panel-wide">
        <p class="panel-kicker">状态快照</p>
        <h2>Go 返回给 Vue 的数据</h2>
        <pre>{{ stateJson }}</pre>
      </article>
    </section>
  </main>
</template>
