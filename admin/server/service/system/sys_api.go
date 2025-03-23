package system

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"gorm.io/gorm"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateApi
//@description: 新增基础api
//@param: api model.SysApi
//@return: err error

type ApiService struct{}

var ApiServiceApp = new(ApiService)

func (apiService *ApiService) CreateApi(api system.SysApi) (err error) {
	if !errors.Is(global.GVA_DB.Where("path = ? AND method = ?", api.Path, api.Method).First(&system.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同api")
	}
	return global.GVA_DB.Create(&api).Error
}

func (apiService *ApiService) CreateAgent(api system.SysAgent) (err error) {
	if !errors.Is(global.GVA_DB.Where("site = ?", api.Site).First(&system.SysAgent{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("已存在相同地区或代理")
	}
	return global.GVA_DB.Create(&api).Error
}

func (apiService *ApiService) GetApiGroups() (groups []string, groupApiMap map[string]string, err error) {
	var apis []system.SysApi
	err = global.GVA_DB.Find(&apis).Error
	if err != nil {
		return
	}
	groupApiMap = make(map[string]string, 0)
	for i := range apis {
		pathArr := strings.Split(apis[i].Path, "/")
		newGroup := true
		for i2 := range groups {
			if groups[i2] == apis[i].ApiGroup {
				newGroup = false
			}
		}
		if newGroup {
			groups = append(groups, apis[i].ApiGroup)
		}
		groupApiMap[pathArr[1]] = apis[i].ApiGroup
	}
	return
}

func (apiService *ApiService) GetAgentList() (apis []system.SysAgent, err error) {
	err = global.GVA_DB.Find(&apis).Error
	if err != nil {
		return
	}
	return
}

func (apiService *ApiService) SyncApi() (newApis, deleteApis, ignoreApis []system.SysApi, err error) {
	newApis = make([]system.SysApi, 0)
	deleteApis = make([]system.SysApi, 0)
	ignoreApis = make([]system.SysApi, 0)
	var apis []system.SysApi
	err = global.GVA_DB.Find(&apis).Error
	if err != nil {
		return
	}
	var ignores []system.SysIgnoreApi
	err = global.GVA_DB.Find(&ignores).Error
	if err != nil {
		return
	}

	for i := range ignores {
		ignoreApis = append(ignoreApis, system.SysApi{
			Path:        ignores[i].Path,
			Description: "",
			ApiGroup:    "",
			Method:      ignores[i].Method,
		})
	}

	var cacheApis []system.SysApi
	for i := range global.GVA_ROUTERS {
		ignoresFlag := false
		for j := range ignores {
			if ignores[j].Path == global.GVA_ROUTERS[i].Path && ignores[j].Method == global.GVA_ROUTERS[i].Method {
				ignoresFlag = true
			}
		}
		if !ignoresFlag {
			cacheApis = append(cacheApis, system.SysApi{
				Path:   global.GVA_ROUTERS[i].Path,
				Method: global.GVA_ROUTERS[i].Method,
			})
		}
	}

	//对比数据库中的api和内存中的api，如果数据库中的api不存在于内存中，则把api放入删除数组，如果内存中的api不存在于数据库中，则把api放入新增数组
	for i := range cacheApis {
		var flag bool
		// 如果存在于内存不存在于api数组中
		for j := range apis {
			if cacheApis[i].Path == apis[j].Path && cacheApis[i].Method == apis[j].Method {
				flag = true
			}
		}
		if !flag {
			newApis = append(newApis, system.SysApi{
				Path:        cacheApis[i].Path,
				Description: "",
				ApiGroup:    "",
				Method:      cacheApis[i].Method,
			})
		}
	}

	for i := range apis {
		var flag bool
		// 如果存在于api数组不存在于内存
		for j := range cacheApis {
			if cacheApis[j].Path == apis[i].Path && cacheApis[j].Method == apis[i].Method {
				flag = true
			}
		}
		if !flag {
			deleteApis = append(deleteApis, apis[i])
		}
	}
	return
}

func (apiService *ApiService) IgnoreApi(ignoreApi system.SysIgnoreApi) (err error) {
	if ignoreApi.Flag {
		return global.GVA_DB.Create(&ignoreApi).Error
	}
	return global.GVA_DB.Unscoped().Delete(&ignoreApi, "path = ? AND method = ?", ignoreApi.Path, ignoreApi.Method).Error
}

func (apiService *ApiService) EnterSyncApi(syncApis systemRes.SysSyncApis) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var txErr error
		if syncApis.NewApis != nil && len(syncApis.NewApis) > 0 {
			txErr = tx.Create(&syncApis.NewApis).Error
			if txErr != nil {
				return txErr
			}
		}
		for i := range syncApis.DeleteApis {
			CasbinServiceApp.ClearCasbin(1, syncApis.DeleteApis[i].Path, syncApis.DeleteApis[i].Method)
			txErr = tx.Delete(&system.SysApi{}, "path = ? AND method = ?", syncApis.DeleteApis[i].Path, syncApis.DeleteApis[i].Method).Error
			if txErr != nil {
				return txErr
			}
		}
		return nil
	})
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteApi
//@description: 删除基础api
//@param: api model.SysApi
//@return: err error

func (apiService *ApiService) DeleteApi(api system.SysApi) (err error) {
	var entity system.SysApi
	err = global.GVA_DB.First(&entity, "id = ?", api.ID).Error // 根据id查询api记录
	if errors.Is(err, gorm.ErrRecordNotFound) {                // api记录不存在
		return err
	}
	err = global.GVA_DB.Delete(&entity).Error
	if err != nil {
		return err
	}
	CasbinServiceApp.ClearCasbin(1, entity.Path, entity.Method)
	return nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetAPIInfoList
//@description: 分页获取数据,
//@param: api model.SysApi, info request.PageInfo, order string, desc bool
//@return: list interface{}, total int64, err error

func (apiService *ApiService) GetAPIInfoList(api system.SysApi, info request.PageInfo, order string, desc bool) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SysApi{})
	var apiList []system.SysApi

	if api.Path != "" {
		db = db.Where("path LIKE ?", "%"+api.Path+"%")
	}

	if api.Description != "" {
		db = db.Where("description LIKE ?", "%"+api.Description+"%")
	}

	if api.Method != "" {
		db = db.Where("method = ?", api.Method)
	}

	if api.ApiGroup != "" {
		db = db.Where("api_group = ?", api.ApiGroup)
	}

	err = db.Count(&total).Error

	if err != nil {
		return apiList, total, err
	}

	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"
	if order != "" {
		orderMap := make(map[string]bool, 5)
		orderMap["id"] = true
		orderMap["path"] = true
		orderMap["api_group"] = true
		orderMap["description"] = true
		orderMap["method"] = true
		if !orderMap[order] {
			err = fmt.Errorf("非法的排序字段: %v", order)
			return apiList, total, err
		}
		OrderStr = order
		if desc {
			OrderStr = order + " desc"
		}
	}
	err = db.Order(OrderStr).Find(&apiList).Error
	return apiList, total, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetAllApis
//@description: 获取所有的api
//@return:  apis []model.SysApi, err error

func (apiService *ApiService) GetAllApis(authorityID uint) (apis []system.SysApi, err error) {
	parentAuthorityID, err := AuthorityServiceApp.GetParentAuthorityID(authorityID)
	if err != nil {
		return nil, err
	}
	err = global.GVA_DB.Order("id desc").Find(&apis).Error
	if parentAuthorityID == 0 || !global.GVA_CONFIG.System.UseStrictAuth {
		return
	}
	paths := CasbinServiceApp.GetPolicyPathByAuthorityId(authorityID)
	// 挑选 apis里面的path和method也在paths里面的api
	var authApis []system.SysApi
	for i := range apis {
		for j := range paths {
			if paths[j].Path == apis[i].Path && paths[j].Method == apis[i].Method {
				authApis = append(authApis, apis[i])
			}
		}
	}
	return authApis, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetApiById
//@description: 根据id获取api
//@param: id float64
//@return: api model.SysApi, err error

func (apiService *ApiService) GetApiById(id int) (api system.SysApi, err error) {
	err = global.GVA_DB.First(&api, "id = ?", id).Error
	return
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateApi
//@description: 根据id更新api
//@param: api model.SysApi
//@return: err error

func (apiService *ApiService) UpdateApi(api system.SysApi) (err error) {
	var oldA system.SysApi
	err = global.GVA_DB.First(&oldA, "id = ?", api.ID).Error
	if oldA.Path != api.Path || oldA.Method != api.Method {
		var duplicateApi system.SysApi
		if ferr := global.GVA_DB.First(&duplicateApi, "path = ? AND method = ?", api.Path, api.Method).Error; ferr != nil {
			if !errors.Is(ferr, gorm.ErrRecordNotFound) {
				return ferr
			}
		} else {
			if duplicateApi.ID != api.ID {
				return errors.New("存在相同api路径")
			}
		}

	}
	if err != nil {
		return err
	}

	err = CasbinServiceApp.UpdateCasbinApi(oldA.Path, api.Path, oldA.Method, api.Method)
	if err != nil {
		return err
	}

	return global.GVA_DB.Save(&api).Error
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteApisByIds
//@description: 删除选中API
//@param: apis []model.SysApi
//@return: err error

func (apiService *ApiService) DeleteApisByIds(ids request.IdsReq) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var apis []system.SysApi
		err = tx.Find(&apis, "id in ?", ids.Ids).Error
		if err != nil {
			return err
		}
		err = tx.Delete(&[]system.SysApi{}, "id in ?", ids.Ids).Error
		if err != nil {
			return err
		}
		for _, sysApi := range apis {
			CasbinServiceApp.ClearCasbin(1, sysApi.Path, sysApi.Method)
		}
		return err
	})
}

func (apiService *ApiService) GetFreeSite(c *gin.Context, regionIdInt int, dateStart, dateEnd, roomId, bedId string) (res []RoomStatusItem, err error) {

	dateStart = dateStart[:10]
	dateEnd = dateEnd[:10]
	roomIdInt, _ := strconv.Atoi(roomId)
	bedIdInt, _ := strconv.Atoi(bedId)

	// 从数据库中获取agent
	var agents []system.SysAgent
	err = global.GVA_DB.Find(&agents).Error
	if err != nil {
		return
	}
	shareMapId := make(map[uint]string, len(agents))
	channelNameMap := make(map[uint]system.SysAgent, len(agents))
	for _, v := range agents {
		shareMapId[v.ID] = v.ShareId
		channelNameMap[v.ID] = v
	}
	errChan := make(chan error)
	defer close(errChan)
	proxyUrl, err := GetProxyIp(c)
	fmt.Println("proxyUrl:", proxyUrl)
	if regionIdInt == 0 {
		resChange := make(chan RoomStatusItem, len(shareMapId))
		var counter int
		var mu sync.Mutex
		for k, v := range shareMapId {
			go func() {
				for i := 0; i < 3; i++ {
					// 更换代理
					if err != nil {
						if err.Error() == "empty" {
							errChan <- errors.New("proxy ip is empty")
							return
						} else {
							continue
						}
					}
					resData, err := HttpCurl(v, channelNameMap[k], dateStart, dateEnd, roomIdInt, bedIdInt, proxyUrl)
					if err != nil {
						continue
					}
					if resData != nil {
						fmt.Println("resData ID:", k)
						for _, itemData := range resData {
							fmt.Println("itemData:", itemData.RegionName+" "+itemData.AgentName)
							resChange <- itemData
						}
						break
					}
				}
				mu.Lock()
				counter++
				mu.Unlock()
				if counter == len(shareMapId) {
					close(resChange)
					return
				}
			}()
		}
		i := 0
		quit := false
		for {
			select {
			case vres, isOk := <-resChange:
				if !isOk {
					quit = true
					break
				}
				i++
				res = append(res, vres)
				fmt.Println("resChange i:", i, "ok")
			case err = <-errChan:
				fmt.Println("err:", err)
				return
			}
			if quit {
				break
			}
		}
	} else {
		for i := 0; i < 3; i++ {
			// 更换代理
			resPos, errSingle := HttpCurl(shareMapId[uint(regionIdInt)], channelNameMap[uint(regionIdInt)],
				dateStart, dateEnd, roomIdInt, bedIdInt, proxyUrl)
			if errSingle != nil {
				fmt.Println("errSingle:", errSingle)
				continue
			}
			for _, itemData := range resPos {
				res = append(res, itemData)
			}
			break
		}
	}
	// 将res按照日期排序
	sort.Slice(res, func(i, j int) bool {
		return res[i].DateStr < res[j].DateStr
	})
	return
}

func GetProxyIp(c *gin.Context) (string, error) {
	proxyIps := []string{}
	//proxyUrl, err := global.GVA_REDIS.Get(c, "proxy_ips").Result()
	//if err != nil || proxyUrl == "" {
	//fmt.Println("proxyUrl:", proxyUrl, "err:", err)
	fullUrl := "http://www.zdopen.com/ShortProxy/GetIP/?api=202410311912222199&akey=34975af8461e8005&timespan=3&type=3"
	client := &http.Client{}

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", nil
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "", nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", nil
	}
	// 解析道ProxyResponse结构体
	var proxyRes ProxyResponse
	err = json.Unmarshal(body, &proxyRes)
	if proxyRes.Code != "10001" {
		fmt.Println("Error unmarshalling response:", err, ", body:", string(body))
		return "", nil
	}
	for _, item := range proxyRes.Data.ProxyList {
		proxyIps = append(proxyIps, "http://"+item.IP+":"+strconv.Itoa(item.Port))
	}
	fmt.Println("proxyIps:", proxyIps)
	//ipListStr := strings.Join(proxyIps, ",")
	//global.GVA_REDIS.Set(c, "proxy_ips", ipListStr, 15*time.Second)
	//} else {
	//	proxyIps = strings.Split(proxyUrl, ",")
	//	if len(proxyIps) == 0 {
	//		return "", errors.New("empty")
	//	}
	//}
	//fmt.Println("after proxyIps:", proxyIps)
	// 随机从0-4中选择一个
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(proxyIps))
	return proxyIps[seededRand], nil
}

func HttpCurl(shareId string, sys system.SysAgent, dateStartStr, dateEndStr string, roomId, bedId int, proxyUrl string) ([]RoomStatusItem, error) {
	payload := []byte("{\"share_id\":\"" + shareId + "\",\"start_date\":\"" + dateStartStr + "\",\"end_date\":\"" + dateEndStr + "\"}")
	var client *http.Client
	if len(proxyUrl) > 0 {
		pUrl, err := url.Parse(proxyUrl)
		if err != nil {
			panic(err)
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(pUrl),
		}
		client = &http.Client{
			Transport: transport,
		}
	} else {
		client = &http.Client{}
	}
	response, err := http.NewRequest("POST", "https://api.yunzhanggui.net/api/v2/hotel_room_share/get_room_status", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("http.NewRequest", err)
		return nil, err
	}
	response.Header.Add("Accept", "application/json, text/plain, */*")
	response.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	response.Header.Add("Connection", "keep-alive")
	response.Header.Add("Content-Type", "application/json;charset=UTF-8")
	response.Header.Add("Origin", "https://uni-pms.yunzhanggui.com")
	response.Header.Add("Referer", "https://uni-pms.yunzhanggui.com/")
	response.Header.Add("Sec-Fetch-Dest", "empty")
	response.Header.Add("Sec-Fetch-Mode", "cors")
	response.Header.Add("Sec-Fetch-Site", "cross-site")
	response.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36")
	response.Header.Add("X-Requested-With", "XMLHttpRequest")
	response.Header.Add("sec-ch-ua", `"Chromium";v="130", "Google Chrome";v="130", "Not?A_Brand";v="99"`)
	response.Header.Add("sec-ch-ua-mobile", "?0")
	response.Header.Add("sec-ch-ua-platform", `"macOS"`)
	response.Header.Add("timestamp", strconv.Itoa(int(time.Now().Unix()))+"000")
	response.Header.Add("ua", "oms_h5")

	var result *http.Response
	defer func() {
		if result != nil {
			result.Body.Close()
		}
	}()
	result, err = client.Do(response)
	if err != nil {
		fmt.Println("==> client.Do", err)
		return nil, err
	}
	if result == nil {
		return nil, errors.New("result is nil")
	}
	body, err := io.ReadAll(result.Body)
	if err != nil {
		fmt.Println("io.ReadAll", err)
		return nil, err
	}
	// 将body字符串转换为Response结构体
	var resData Response
	err = json.Unmarshal(body, &resData)
	if err != nil {
		var resDataEmpty ResponseEmpty
		err = json.Unmarshal(body, &resDataEmpty)
		if err != nil {
			fmt.Println("json.Unmarshal", err, ", body:", string(body))
			return nil, err
		}
		return nil, nil
	}
	// resData结构体中，ret_data.room_status是一个map，key是房间类型的id，value是房间类型的信息, 我们需要根据req请求结构体中参数， 获取在哪个地区， 在哪一天， 有哪些空闲的房间， 并返回room_type_name， 形成一个数组
	// 遍历resData.RetData.RoomStatus， 并获取指定日期的room_count不为0房间信息
	roomStatusItemList := []RoomStatusItem{}
	for _, roomStatus := range resData.RetData.RoomStatus {
		// 获取roomName
		for _, item := range roomStatus.RoomList {
			// 获取指定日期的room_count不为0房间信息
			for date, it := range item.ScreenList {
				if date >= dateStartStr && date <= dateEndStr && it.HasOrderState == 2 {
					matchStr := ""
					if roomId > 0 && bedId > 0 {
						matchStr = strconv.Itoa(int(roomId)) + "房" + strconv.Itoa(int(bedId)) + "床"
					} else if roomId > 0 {
						matchStr = strconv.Itoa(int(roomId)) + "房"
					} else if bedId > 0 {
						matchStr = strconv.Itoa(int(bedId)) + "床"
					}
					if len(matchStr) > 0 {
						re := regexp.MustCompile(matchStr)
						match := re.FindString(item.RoomName)
						fmt.Println("match:", match, "item.RoomName:", item.RoomName, "matchStr:", matchStr)
						if len(match) > 0 {
							roomStatusItemList = append(roomStatusItemList, RoomStatusItem{
								RegionName:   sys.Region,
								AgentName:    sys.Agent,
								DateStr:      date,
								RoomTypeID:   roomStatus.RoomTypeID,
								RoomTypeName: roomStatus.RoomTypeName,
								RoomName:     item.RoomName,
								RoomCount:    roomStatus.RoomCount,
							})
						}
					} else {
						roomStatusItemList = append(roomStatusItemList, RoomStatusItem{
							RegionName:   sys.Region,
							AgentName:    sys.Agent,
							DateStr:      date,
							RoomTypeID:   roomStatus.RoomTypeID,
							RoomTypeName: roomStatus.RoomTypeName,
							RoomName:     item.RoomName,
							RoomCount:    roomStatus.RoomCount,
						})
					}
				}
			}
		}

	}
	return roomStatusItemList, nil
}

type Response struct {
	RetCode int     `json:"ret_code"`
	RetMsg  string  `json:"ret_msg"`
	RetData RetData `json:"ret_data"`
}

type ResponseEmpty struct {
	RetCode int      `json:"ret_code"`
	RetMsg  string   `json:"ret_msg"`
	RetData []string `json:"ret_data"`
}

type RetData struct {
	BookRate          float64               `json:"book_rate"`
	BookRoomCount     float64               `json:"book_room_count"`
	BookRoomFee       float64               `json:"book_room_fee"`
	InnName           string                `json:"inn_name"`
	DateRangeMaxLimit DateRange             `json:"date_range_max_limit"`
	CurrentDate       DateRange             `json:"current_date"`
	RoomStatus        map[string]RoomStatus `json:"room_status"`
	Setting           Setting               `json:"setting"`
}

type DateRange struct {
	MinStartDate string `json:"min_start_date"`
	MaxEndDate   string `json:"max_end_date"`
	StartDate    string `json:"start_date,omitempty"`
	EndDate      string `json:"end_date,omitempty"`
}

type RoomStatus struct {
	RoomTypeID   int               `json:"room_type_id"`
	RoomTypeName string            `json:"room_type_name"`
	RoomList     map[string]Room   `json:"room_list"`
	RoomCount    int               `json:"room_count"`
	ScreenList   map[string]Screen `json:"screen_list"`
}

type RoomStatusItem struct {
	RegionName   string `json:"region_name"`
	AgentName    string `json:"agent_name"`
	DateStr      string `json:"date_str"`
	RoomTypeID   int    `json:"room_type_id"`
	RoomTypeName string `json:"room_type_name"`
	RoomName     string `json:"room_name"`
	RoomCount    int    `json:"room_count"`
}

type Room struct {
	RoomID      int                   `json:"room_id"`
	RoomName    string                `json:"room_name"`
	CleanStatus int                   `json:"clean_status"`
	ScreenList  map[string]RoomScreen `json:"screen_list"`
	IsShow      bool                  `json:"is_show"`
}

type RoomScreen struct {
	RoomOrderState string `json:"room_order_state"`
	HasOrderState  int    `json:"has_order_state"`
	RoomColor      string `json:"room_color"`
	Price          string `json:"price"`
	RoomFee        string `json:"room_fee"`
	Date           string `json:"date"`
}

type Screen struct {
	RoomNum   int    `json:"room_num"`
	RoomPrice string `json:"room_price"`
}

type Setting struct {
	ID               int       `json:"id"`
	InnID            int       `json:"innid"`
	Title            string    `json:"title"`
	Tag              string    `json:"tag"`
	DateType         string    `json:"date_type"`
	DateValue        DateValue `json:"date_value"`
	RoomIDs          string    `json:"room_ids"`
	RoomFeeState     string    `json:"room_fee_state"`
	DatePriceRpID    int       `json:"date_price_rp_id"`
	OrderStatusState string    `json:"order_status_state"`
	DirtyRoomState   string    `json:"dirty_room_state"`
	CreatedAt        string    `json:"created_at"`
	UpdatedAt        string    `json:"updated_at"`
	ShareID          string    `json:"share_id"`
	ShareState       string    `json:"share_state"`
	ShareType        string    `json:"share_type"`
	RoomState        string    `json:"room_state"`
	ShareRoomShow    string    `json:"share_room_show"`
	DatePriceRpState string    `json:"date_price_rp_state"`
	IsAllRoom        string    `json:"is_all_room"`
}

type DateValue struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Before    string `json:"before"`
	After     string `json:"after"`
}

type Proxy struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type ProxyListData struct {
	Count     int     `json:"count"`
	ProxyList []Proxy `json:"proxy_list"`
}

type ProxyResponse struct {
	Code string        `json:"code"`
	Msg  string        `json:"msg"`
	Data ProxyListData `json:"data"`
}
