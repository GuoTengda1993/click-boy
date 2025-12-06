<template>
<main style="padding: 10px; margin-top: -10px;">
  <div style="display: flex; gap: 16px; margin-top: 8px;">
    <div style="flex: 1;">
      <el-row style="margin-bottom: 5px;">
        <el-col :span="18">
          <div>
            <el-button type="primary" link @click="addPoints(-1, 3, 'Key')" :disabled="!devConnected || clickStatus != 0">Home</el-button>
            <el-button type="primary" link @click="addPoints(-1, 25, 'Key')" :disabled="!devConnected || clickStatus != 0">Vol-</el-button>
            <el-button type="primary" link @click="addPoints(-1, 24, 'Key')" :disabled="!devConnected || clickStatus != 0">Vol+</el-button>
            <el-button type="primary" link @click="addPoints(-9, 0, 'StartApp')" :disabled="!devConnected || clickStatus != 0 || packageName == ''">StartApp</el-button>
            <el-button type="primary" link @click="addPoints(-9, 0, 'StopApp')" :disabled="!devConnected || clickStatus != 0 || packageName == ''">StopApp</el-button>
          </div>
          <div>
            <el-button type="primary" link @click="leftSwipeBack" :disabled="!devConnected || clickStatus != 0">Swipe Back</el-button>
            <el-button type="primary" link @click="upSwipeHome" :disabled="!devConnected || clickStatus != 0">Swipe Home</el-button>
            <el-button type="primary" link @click="upSwipeTask" :disabled="!devConnected || clickStatus != 0">Swipe Task</el-button>
          </div>
        </el-col>
        <el-col :span="6">
          <el-button v-if="devConnected" type="primary" @click="doScreencap">Screen Shot</el-button>
          <el-button v-else type="primary" @click="ConnectDevice">Connect</el-button>
        </el-col>
      </el-row>

      <div v-if="imgSrc" style="border: 1px solid #ddd; max-height: 70vh; overflow: auto;">
        <img :src="imgSrc"
             ref="imgRef"
             style="max-width: auto; height: 777px; cursor: crosshair;"
             @click="onImgClick" />
      </div>
      <div v-else>
        <div v-if="devConnected" style="color: #666;">Please click "Screen Shot" to get image.</div>
        <div v-else style="color: #666;">
          <el-text>
            <br/>
            <el-text size="large">Android</el-text><br/>
            1. install ADB tools;<br/>
            2. run command 'adb start-server'.<br/>
            3. switch on developer mode in phone.
          </el-text>
          <el-text>
          <br/><br/>
          <el-text size="large">iOS</el-text><br/>
            1. install go-ios tools: npm install -g go-ios;<br/>
            2. run command 'sudo ios tunnel start';<br/>
            3. run command 'ios runwda' in anothor terminal.
          </el-text>
        </div>
      </div>
      
    </div>

    <div style="width: 550px;">
      <div style="display:flex; gap:8px; margin-bottom: 8px;">
        <el-input v-model.number="intervalMs" type="number" style="width: 270px;">
            <template #prepend>Click Interval</template>
            <template #append>ms</template>
        </el-input>
        <el-input v-model.number="times" type="number" style="width: 270px;">
            <template #prepend>Loop Times</template>
        </el-input>
      </div>
      <div style="display:flex; gap:8px; margin-bottom: 8px;">
        <el-input v-model.number="durationMs" type="number" >
            <template #prepend>Swip/Press Duration(andriod)</template>
            <template #append>ms</template>
        </el-input>
      </div>
      <div style="display:flex; gap:8px; margin-bottom: 8px;">
        <el-input v-model.number="packageName" style="width: 270px;">
          <template #prepend>Package/BundleId</template>
        </el-input>
        <el-input v-model.number="appActivity" style="width: 270px;">
          <template #prepend>Activity(android)</template>
        </el-input>
      </div>
      <div style="display:flex; gap:8px; margin-bottom: 8px;">
        <el-button type="success" @click="runClick" :disabled="!canRun">Run</el-button>
        <el-button v-if="clickStatus == 1" type="warning" @click="clickPause">Pause</el-button>
        <el-button v-if="clickStatus == 2" type="primary" @click="clickResume">Resume</el-button>
        <el-button v-if="clickStatus != 0" type="danger" @click="clickStop">Stop</el-button>
        <div style="display:flex; gap:8px;" v-if="loopTime > 0">
          <el-text>The
            <el-text size="large" style="color: #2768F5;">{{ loopTime }}</el-text>th loop.
          </el-text>
        </div>
      </div>
      <div style="border: 1px dashed #bbb; padding: 8px; height: 60vh; overflow: auto;">
        <el-row type="flex" :gutter="10" style="display: flex; align-items: center; margin-bottom: 8px;">
          <div style="color: #666;">Click on the screenshot to record the coordinates</div>
          <div style="flex-grow: 1; display: flex; justify-content: right; margin-right: 5px;">
            <el-button link type="warning" @click="clearPoints" :disabled="points.length===0">Clear</el-button>
          </div>
        </el-row>
        <el-row type="flex" :gutter="10" style="display: flex; align-items: center; margin-bottom: 8px;">
          <div style="color: #777; font-size: smaller;">The "Swipe" event requires that both adjacent coordinate points be set as swiping. When triggered, it will swipe from the first coordinate point to the second.</div>
          <div style="color: #777; font-size: smaller;">[Android only] The duration of the "QuickSwipe" event is fixed at 100ms, while that of the "SlowSwipe" event is fixed at 1000ms.</div>
        </el-row>
        
        <el-table :data="points" size="small" style="width: 100%;">
          <el-table-column prop="x" label="X" />
          <el-table-column prop="y" label="Y" />
          <el-table-column prop="event" label="Event">
            <template #default="{ row }">
              <span v-if="!row.edit">{{ row.event }}</span>
              <el-select v-else v-model="row.event" placeholder="Select event" style="width: 100px" size="small">
                <el-option
                  v-for="item in eventOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column label="Operate">
            <template #default="scope">
              <el-button link type="danger" @click="removePoint(scope.$index)">Delete</el-button>
              <template v-if="scope.row.x >= 0">
                <el-button v-if="!scope.row.edit" link type="primary" @click="handleEditEvent(scope.row)">Modify</el-button>
                <el-button v-else link type="primary" @click="handleCancelEdit(scope.row)">Cancel</el-button>
              </template>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </div>
</main>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeUnmount, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { ScreenShot, StartClick, ConnectDevice, Pause, Resume, Stop } from '../../wailsjs/go/main/App'
import { click } from '../../wailsjs/go/models'
import {
    EventsOn,
    EventsOffAll
} from '../../wailsjs/runtime/runtime'

const devConnected = ref<boolean>(false)
const devType = ref<string>('')

const imgSrc = ref<string>('')
const imgRef = ref<HTMLImageElement | null>(null)

const intervalMs = ref<number>(500)
const times = ref<number>(10)
const packageName = ref<string>('')
const appActivity = ref<string>('')
const points = reactive<Array<click.Point>>([])

const loopTime = ref<number>(0)

const devWidth = ref<number>(0)
const devHeight = ref<number>(0)

const clickStatus = ref<number>(0)

const eventOptionsIos = [
  {
    value: 'Click',
    label: 'Click',
  },
  {
    value: 'DbClick',
    label: 'DbClick',
  },
  {
    value: 'LongPress',
    label: 'LongPress',
  },
  {
    value: 'Swipe',
    label: 'Swipe',
  },
  {
    value: 'Drag',
    label: 'Drag',
  },
]

const eventOptionsAndroid = [
  {
    value: 'Click',
    label: 'Click',
  },
  {
    value: 'LongPress',
    label: 'LongPress',
  },
  {
    value: 'DbClick',
    label: 'DbClick',
  },
  {
    value: 'Swipe',
    label: 'Swipe',
  },
  {
    value: 'QuickSwipe',
    label: 'QuickSwipe',
  },
  {
    value: 'SlowSwipe',
    label: 'SlowSwipe',
  },
]

const eventOptions = computed(() => {
  if (devType.value == 'ios') {
    return eventOptionsIos
  }
  return eventOptionsAndroid
})

const canRun = computed(() => devConnected.value && points.length > 0 && intervalMs.value >= 0 && times.value > 0 && clickStatus.value == 0)

onMounted(() => {
  EventsOn("message", (msgInfo) => {
    ElMessage({
      message: msgInfo['msg'],
      type: msgInfo['type'],
    })
  })

  EventsOn("loop-num", (n) => {
    loopTime.value = n
    if (n == -1) {
      clickStatus.value = 0
    }
  })

  EventsOn("detect-device", (device) => {
    devType.value = device['type']
    devConnected.value = device['status']
    devWidth.value = device['width']
    devHeight.value = device['height']
  })
})

onBeforeUnmount(() => {
  EventsOffAll();
})

async function doScreencap() {
  if (!devConnected.value) return
  const dataUrl = await ScreenShot()
  if (!dataUrl) {
    return
  }
  imgSrc.value = dataUrl
}

function onImgClick(e: MouseEvent) {
  if (clickStatus.value != 0) {
    ElMessage.warning('Please stop first and then add the click coordinates.')
    return
  }
  const img = imgRef.value
  if (!img) return
  const displayW = img.clientWidth
  const displayH = img.clientHeight
  let natW = img.naturalWidth
  let natH = img.naturalHeight
  if (devType.value == 'ios') {
    natW = devWidth.value
    natH = devHeight.value 
  }

  const relX = (e as any).offsetX as number
  const relY = (e as any).offsetY as number

  const absX = Math.round(relX * (natW / displayW))
  const absY = Math.round(relY * (natH / displayH))

  cancelAllEdit()
  points.push({ x: absX, y: absY, event: 'Click', edit: false } as click.Point)
}

function addPoints(x: number, y: number, event: string) {
  cancelAllEdit()
  points.push({ x: x, y: y, event: event, edit: false } as click.Point)
}

function removePoint(idx: number) {
  points.splice(idx, 1)
}

function clearPoints() {
  points.splice(0, points.length)
}

const leftSwipeBack = () => {
  cancelAllEdit()
  const img = imgRef.value
  if (!img) return
  let natH = devHeight.value
  let event = 'Drag'
  if (devType.value == 'android') {
    natH = img.naturalHeight
    event = 'QuickSwipe'
  }

  const y = Math.round(natH / 2)
  points.push({ x: 3, y: y, event: event, edit: false } as click.Point)
  points.push({ x: 400, y: y, event: event, edit: false } as click.Point)
}

const upSwipeHome = () => {
  cancelAllEdit()
  const img = imgRef.value
  if (!img) return
  let natW = devWidth.value
  let natH = devHeight.value
  let event = 'Drag'
  if (devType.value == 'android') {
    natW = img.naturalWidth
    natH = img.naturalHeight
    event = 'QuickSwipe'
  }
  const x = Math.round(natW / 2)
  points.push({ x: x, y: natH, event: event, edit: false } as click.Point)
  points.push({ x: x, y: natH - 100, event: event, edit: false } as click.Point)
}

const upSwipeTask = () => {
  cancelAllEdit()
  const img = imgRef.value
  if (!img) return
  let natW = devWidth.value
  let natH = devHeight.value
  let event = 'Drag'
  if (devType.value == 'android') {
    natW = img.naturalWidth
    natH = img.naturalHeight
    event = 'SlowSwipe'
  }
  const x = Math.round(natW / 2)
  points.push({ x: x, y: natH, event: event, edit: false } as click.Point)
  points.push({ x: x, y: natH - 400, event: event, edit: false } as click.Point)
}

const handleEditEvent = (row: click.Point) => {
  // 可以在这里先关闭其他行的编辑状态
  points.forEach(item => {
    if (item !== row) {
      item.edit = false;
    }
  });
  // 开启当前行的编辑状态
  row.edit = true;
};

const handleCancelEdit = (row: click.Point) => {
  row.edit = false;
};

const cancelAllEdit = () => {
  points.forEach(item => {
    item.edit = false;
  });
};

async function runClick() {
  if (!devConnected.value) return
  if (points.length === 0) {
    ElMessage.warning('no valid points')
    return
  }
  cancelAllEdit()
  let params = {
    times: times.value,
    interval: intervalMs.value,
    duration: 1,
    package: packageName.value,
    activity: appActivity.value,
  }
  StartClick(points, params)
  clickStatus.value = 1
}

const clickPause = () => {
  Pause().then(r => {
    clickStatus.value = r
    ElMessage.info('Will be pause after this loop')
  })
}

const clickResume = () => {
  Resume().then(r => {
    clickStatus.value = r
    ElMessage.info('Resume success')
  })
}

const clickStop = () => {
  Stop().then(r => {
    clickStatus.value = r
    ElMessage.info('Will be stop after this loop')
  })
}
</script>

<style scoped>
</style>