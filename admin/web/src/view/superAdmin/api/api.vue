<template>
  <div>
    <div class="gva-search-box">
      <el-form
        ref="searchForm"
        :inline="true"
        :model="searchInfo"
      >
        <el-form-item label="地区">
          <template #default="{row}">
              <el-select
                v-model="searchInfo.regionId"
                placeholder="请选择或新增"
                allow-create filterable default-first-option
              >
                <el-option
                  v-for="item in apiAgentsOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                />
              </el-select>
            </template>
        </el-form-item>
        <el-form-item label="日期">
          <el-date-picker v-model="searchInfo.startCreatedAt" type="date" placeholder="开始日期" format="YYYY-MM-DD"></el-date-picker>
          —
          <el-date-picker v-model="searchInfo.endCreatedAt" type="date" placeholder="结束日期" format="YYYY-MM-DD"></el-date-picker>
        </el-form-item>
        <el-form-item label="房型">
          <el-select
            v-model="searchInfo.roomId"
            clearable
            placeholder="请选择"
          >
            <el-option key="1" label="1房" value="1"/>
            <el-option key="2" label="2房" value="2"/>
            <el-option key="3" label="3房" value="3"/>
            <el-option key="4" label="4房" value="4"/>
            <el-option key="5" label="5房" value="5"/>
            <el-option key="6" label="6房" value="6"/>
            <el-option key="7" label="7房" value="7"/>
            <el-option key="8" label="8房" value="8"/>
            <el-option key="9" label="9房" value="9"/>
            <el-option key="10" label="10房" value="10"/>
            <el-option key="11" label="11房" value="11"/>
            <el-option key="12" label="12房" value="12"/>
          </el-select>
        </el-form-item>
        <el-form-item label="床位">
          <el-select
            v-model="searchInfo.bedId"
            clearable
            placeholder="请选择"
          >
            <el-option key="1" label="1床" value="1"/>
            <el-option key="2" label="2床" value="2"/>
            <el-option key="3" label="3床" value="3"/>
            <el-option key="4" label="4床" value="4"/>
            <el-option key="5" label="5床" value="5"/>
            <el-option key="6" label="6床" value="6"/>
            <el-option key="7" label="7床" value="7"/>
            <el-option key="8" label="8床" value="8"/>
            <el-option key="9" label="9床" value="9"/>
            <el-option key="10" label="10床" value="10"/>
            <el-option key="11" label="11床" value="11"/>
            <el-option key="12" label="12床" value="12"/>
            <el-option key="13" label="13床" value="13"/>
            <el-option key="14" label="14床" value="14"/>
            <el-option key="15" label="15床" value="15"/>
            <el-option key="16" label="16床" value="16"/>
            <el-option key="17" label="17床" value="17"/>
            <el-option key="18" label="18床" value="18"/>
            <el-option key="19" label="19床" value="19"/>
            <el-option key="20" label="20床" value="20"/>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            icon="search"
            @click="onSubmit"
          >
            查询
          </el-button>
          <el-button
            icon="refresh"
            @click="onReset"
          >
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
    <div class="gva-btn-list">
            <el-button
              type="primary"
              icon="plus"
              @click="openDialog('addApi')"
            >
              新增
            </el-button>
          </div>
      <el-table
        :data="tableData"
        @sort-change="sortChange"
        @selection-change="handleSelectionChange"
      >
        <el-table-column
          type="selection"
          width="55"
        />
        <el-table-column
          align="left"
          label="地区名称"
          min-width="100"
          prop="region_name"
          sortable="custom"
        />
        <el-table-column
          align="left"
          label="代理人名称"
          min-width="100"
          prop="agent_name"
          sortable="custom"
        />
        <el-table-column
          align="left"
          label="日期"
          min-width="100"
          prop="date_str"
          sortable="custom"
        />
        <el-table-column
          align="left"
          label="房间类型ID"
          min-width="150"
          prop="room_type_id"
          sortable="custom"
        />
        <el-table-column
          align="left"
          label="房间类型名称"
          min-width="150"
          prop="room_type_name"
          sortable="custom"
        />
        <el-table-column
          align="left"
          label="房间名称"
          min-width="150"
          prop="room_name"
          sortable="custom"
        />
        <el-table-column
          align="left"
          label="房间数量"
          min-width="150"
          prop="room_count"
          sortable="custom"
        />
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-drawer
      v-model="syncApiFlag"
      size="80%"
      :before-close="closeSyncDialog"
      :show-close="false"
    >
      <warning-bar title="同步API，不输入路由分组将不会被自动同步" />
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">同步路由</span>
          <div>
            <el-button :loading="apiCompletionLoading" @click="closeSyncDialog">
              取 消
            </el-button>
            <el-button
              type="primary"
              :loading="syncing||apiCompletionLoading"
              @click="enterSyncDialog"
            >
              确 定
            </el-button>
          </div>
        </div>
      </template>

      <h4>新增路由 <span class="text-xs text-gray-500 mx-2 font-normal">存在于当前路由中，但是不存在于api表</span>
        <el-button type="primary" size="small" @click="apiCompletion">
          <el-icon size="18">
            <ai-gva />
          </el-icon>
          自动填充
        </el-button>
      </h4>
      <el-table
         v-loading="syncing||apiCompletionLoading"
         element-loading-text="小淼正在思考..."
        :data="syncApiData.newApis"
      >
        <el-table-column
          align="left"
          label="API路径"
          min-width="150"
          prop="path"
        />
        <el-table-column
          align="left"
          label="API分组"
          min-width="150"
          prop="apiGroup"
        >
          <template #default="{row}">
            <el-select
              v-model="row.apiGroup"
              placeholder="请选择或新增"
              allow-create filterable default-first-option
            >
              <el-option
                v-for="item in apiGroupOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="API简介"
          min-width="150"
          prop="description"
        >
          <template #default="{row}">
            <el-input
              v-model="row.description"
              autocomplete="off"
            />
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="请求"
          min-width="150"
          prop="method"
        >
          <template #default="scope">
            <div>
              {{ scope.row.method }} / {{ methodFilter(scope.row.method) }}
            </div>
          </template>
        </el-table-column>
        <el-table-column
          label="操作"
          min-width="150"
          fixed="right"
        >
          <template #default="{row}">
            <el-button icon="plus" type="primary" link @click="addApiFunc(row)">
              单条新增
            </el-button>
            <el-button icon="sunrise" type="primary" link @click="ignoreApiFunc(row,true)">
              忽略
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <h4>已删除路由 <span class="text-xs text-gray-500 ml-2 font-normal">已经不存在于当前项目的路由中，确定同步后会自动从apis表删除</span></h4>
      <el-table
        :data="syncApiData.deleteApis"
      >
        <el-table-column
          align="left"
          label="API路径"
          min-width="150"
          prop="path"
        />
        <el-table-column
          align="left"
          label="API分组"
          min-width="150"
          prop="apiGroup"
        />
        <el-table-column
          align="left"
          label="API简介"
          min-width="150"
          prop="description"
        />
        <el-table-column
          align="left"
          label="请求"
          min-width="150"
          prop="method"
        >
          <template #default="scope">
            <div>
              {{ scope.row.method }} / {{ methodFilter(scope.row.method) }}
            </div>
          </template>
        </el-table-column>
      </el-table>

      <h4>忽略路由 <span class="text-xs text-gray-500 ml-2 font-normal">忽略路由不参与api同步，常见为不需要进行鉴权行为的路由</span></h4>
      <el-table
        :data="syncApiData.ignoreApis"
      >
        <el-table-column
          align="left"
          label="API路径"
          min-width="150"
          prop="path"
        />
        <el-table-column
          align="left"
          label="API分组"
          min-width="150"
          prop="apiGroup"
        />
        <el-table-column
          align="left"
          label="API简介"
          min-width="150"
          prop="description"
        />
        <el-table-column
          align="left"
          label="请求"
          min-width="150"
          prop="method"
        >
          <template #default="scope">
            <div>
              {{ scope.row.method }} / {{ methodFilter(scope.row.method) }}
            </div>
          </template>
        </el-table-column>
        <el-table-column
          label="操作"
          min-width="150"
          fixed="right"
        >
          <template #default="{row}">
            <el-button icon="sunny" type="primary" link @click="ignoreApiFunc(row,false)">
              取消忽略
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>

    <el-drawer
      v-model="dialogFormVisible"
      size="60%"
      :before-close="closeDialog"
      :show-close="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ dialogTitle }}</span>
          <div>
            <el-button @click="closeDialog">
              取 消
            </el-button>
            <el-button
              type="primary"
              @click="enterDialog"
            >
              确 定
            </el-button>
          </div>
        </div>
      </template>

      <el-form
        ref="apiForm"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item
          label="地区网址"
          prop="site"
        >
          <el-input
            v-model="form.site"
            autocomplete="off"
          />
        </el-form-item>
        <el-form-item
          label="地区名称"
          prop="region"
        >
          <el-input
            v-model="form.region"
            autocomplete="off"
          />
        </el-form-item>
        <el-form-item
              label="代理人"
              prop="agent"
            >
              <el-input
                v-model="form.agent"
                autocomplete="off"
              />
            </el-form-item>
        <el-form-item
          label="地区简介"
          prop="description"
        >
          <el-input
            v-model="form.description"
            autocomplete="off"
          />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  getApiById,
  getApiList,
  getFreeSite,
  createAgent,
  getAgentList,
  createApi,
  updateApi,
  deleteApi,
  deleteApisByIds,
  freshCasbin,
  syncApi,
  getApiGroups,
  ignoreApi,
  enterSyncApi
} from '@/api/api'
import { toSQLLine } from '@/utils/stringFun'
import WarningBar from '@/components/warningBar/warningBar.vue'
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import ExportExcel from '@/components/exportExcel/exportExcel.vue'
import ExportTemplate from '@/components/exportExcel/exportTemplate.vue'
import ImportExcel from '@/components/exportExcel/importExcel.vue'
import {butler} from "@/api/autoCode";
import moment from 'moment';

defineOptions({
  name: 'Api',
})

const methodFilter = (value) => {
  const target = methodOptions.value.filter(item => item.value === value)[0]
  return target && `${target.label}`
}

const apis = ref([])
const form = ref({
  path: '',
  apiGroup: '',
  method: '',
  description: ''
})
const methodOptions = ref([
  {
    value: 'POST',
    label: '创建',
    type: 'success'
  },
  {
    value: 'GET',
    label: '查看',
    type: ''
  },
  {
    value: 'PUT',
    label: '更新',
    type: 'warning'
  },
  {
    value: 'DELETE',
    label: '删除',
    type: 'danger'
  }
])

const type = ref('')
const rules = ref({
  site: [{ required: true, message: '请输入网址', trigger: 'blur' }],
  region: [{ required: true, message: '请输入地区名称', trigger: 'blur' }],
  agent: [{ required: true, message: '请输入代理人', trigger: 'blur' }],
  description: [
    { required: false, message: '请输入地区介绍', trigger: 'blur' }
  ]
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const apiGroupOptions = ref([])
const apiGroupMap = ref({})
const apiAgentsOptions = ref([])

const getGroup = async() => {
  const res = await getApiGroups()
  if (res.code === 0) {
    const groups = res.data.groups
    apiGroupOptions.value = groups.map(item => ({ label: item, value: item }))
    apiGroupMap.value = res.data.apiGroupMap
  }
}

const getAgents = async() => {
  const res = await getAgentList()
  if (res.code === 0) {
    const agents = res.data
    apiAgentsOptions.value = agents.map(item => ({ label: item.region + '-'  + item.agent, value: item.ID }))
    console.log(apiAgentsOptions.value)
  }
}

const ignoreApiFunc = async (row,flag) =>{
  const res = await ignoreApi({path:row.path,method:row.method,flag})
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: res.msg
    })
    if(flag){
      syncApiData.value.newApis = syncApiData.value.newApis.filter(item => !(item.path === row.path && item.method === row.method))
      syncApiData.value.ignoreApis.push(row)
      return
    }
    syncApiData.value.ignoreApis = syncApiData.value.ignoreApis.filter(item => !(item.path === row.path && item.method === row.method))
    syncApiData.value.newApis.push(row)
  }
}

const addApiFunc = async(row)=>{
  if(!row.apiGroup){
    ElMessage({
      type: 'error',
      message: '请先选择API分组'
    })
    return
  }
  if(!row.description){
    ElMessage({
      type: 'error',
      message: '请先填写API描述'
    })
    return
  }
  const res = await createApi(row)
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '添加成功',
      showClose: true
    })
    syncApiData.value.newApis = syncApiData.value.newApis.filter(item => !(item.path === row.path && item.method === row.method))
  }
  getTableData()
  getGroup()
}

const closeSyncDialog = () => {
  syncApiFlag.value = false
}

const syncing = ref(false)

const enterSyncDialog = async() => {
 if( syncApiData.value.newApis.some(item => !item.apiGroup || !item.description)){
   ElMessage({
     type: 'error',
     message: '存在API未分组或未填写描述'
   })
   return
 }

  syncing.value = true
  const res = await enterSyncApi(syncApiData.value)
  syncing.value = false
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: res.msg
    })
    syncApiFlag.value = false
    getTableData()
  }
}

const onReset = () => {
  searchInfo.value = {}
  getTableData()
}
// 搜索

const onSubmit = () => {
  page.value = 1
  pageSize.value = 10
  getTableData()
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 排序
const sortChange = ({ prop, order }) => {
  if (prop) {
    if (prop === 'ID') {
      prop = 'id'
    }
    searchInfo.value.orderKey = toSQLLine(prop)
    searchInfo.value.desc = order === 'descending'
  }
  getTableData()
}

// 查询
const getTableData = async() => {
searchInfo.value.startCreatedAt = moment(searchInfo.value.startCreatedAt).format('YYYY-MM-DD HH:mm:ss');
  searchInfo.value.endCreatedAt = moment(searchInfo.value.endCreatedAt).format('YYYY-MM-DD HH:mm:ss');
  const table = await getFreeSite({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

// getTableData()
getAgents()
// 批量操作
const handleSelectionChange = (val) => {
  apis.value = val
}

const onDelete = async() => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
    const ids = apis.value.map(item => item.ID)
    const res = await deleteApisByIds({ ids })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: res.msg
      })
      if (tableData.value.length === ids.length && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}
const onFresh = async() => {
  ElMessageBox.confirm('确定要刷新缓存吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
    const res = await freshCasbin()
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: res.msg
      })
    }
  })
}

const syncApiData = ref({
  newApis:[],
  deleteApis:[],
  ignoreApis:[]
})

const syncApiFlag = ref(false)

const onSync = async() => {
  const res = await syncApi()
  if (res.code === 0) {
    res.data.newApis.forEach(item => {
      item.apiGroup = apiGroupMap.value[item.path.split('/')[1]]
    })

    syncApiData.value = res.data
    syncApiFlag.value = true
  }
}

// 弹窗相关
const apiForm = ref(null)
const initForm = () => {
  apiForm.value.resetFields()
  form.value = {
    path: '',
    apiGroup: '',
    method: '',
    description: ''
  }
}

const dialogTitle = ref('新增地区/代理人')
const dialogFormVisible = ref(false)
const openDialog = (key) => {
  switch (key) {
    case 'addApi':
      dialogTitle.value = '新增地区/代理人'
      break
    case 'edit':
      dialogTitle.value = '编辑Api'
      break
    default:
      break
  }
  type.value = key
  dialogFormVisible.value = true
}
const closeDialog = () => {
  initForm()
  dialogFormVisible.value = false
}

const editApiFunc = async(row) => {
  const res = await getApiById({ id: row.ID })
  form.value = res.data.api
  openDialog('edit')
}

const enterDialog = async() => {
  apiForm.value.validate(async valid => {
    if (valid) {
      switch (type.value) {
        case 'addApi':
          {
            const res = await createAgent(form.value)
            if (res.code === 0) {
              ElMessage({
                type: 'success',
                message: '添加成功',
                showClose: true
              })
            }
            getTableData()
            getAgents()
            closeDialog()
          }

          break
        case 'edit':
          {
            const res = await updateApi(form.value)
            if (res.code === 0) {
              ElMessage({
                type: 'success',
                message: '编辑成功',
                showClose: true
              })
            }
            getTableData()
            closeDialog()
          }
          break
        default:
          // eslint-disable-next-line no-lone-blocks
          {
            ElMessage({
              type: 'error',
              message: '未知操作',
              showClose: true
            })
          }
          break
      }
    }
  })
}

const deleteApiFunc = async(row) => {
  ElMessageBox.confirm('此操作将永久删除所有角色下该api, 是否继续?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(async() => {
      const res = await deleteApi(row)
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功!'
        })
        if (tableData.value.length === 1 && page.value > 1) {
          page.value--
        }
        getTableData()
        getGroup()
      }
    })
}
const apiCompletionLoading = ref(false)
const apiCompletion = async () =>{
  apiCompletionLoading.value = true
  const routerPaths = syncApiData.value.newApis.filter(item => !item.apiGroup || !item.description).map(item => item.path)
  const res = await butler({data:routerPaths,command:'apiCompletion'})
  apiCompletionLoading.value = false
  if (res.code === 0) {
    try{
      const data = JSON.parse(res.data)
      syncApiData.value.newApis.forEach(item => {
        const target = data.find(d => d.path === item.path)
        if(target){
          if(!item.apiGroup){
            item.apiGroup = target.apiGroup
          }
          if (!item.description) {
            item.description = target.description
          }
        }
      })
    } catch (e) {
      ElMessage({
        type: 'error',
        message: 'AI自动填充失败,请重新生成'
      })
    }
  }
}

</script>

<style scoped lang="scss">
.warning {
  color: #dc143c;
}
</style>
