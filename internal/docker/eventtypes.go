package docker

type Type string

var EventTypes []string = []string{
	"container",
	"daemon",
	"image",
	"network",
	"volume",
}

var ContainerEvent []string = []string{
	"create",
	"start",
	"restart",
	"stop",
	"die",
	"kill",
	"pause",
	"unpause",
	"destroy",
	"rename",
	"update",
	"attach",
	"detach",
	"exec_create",
	"exec_start",
	"exec_die",
}

var ImagerEvent []string = []string{
	"pull",
	"push",
	"tag",
	"untag",
	"delete",
	"save",
	"load",
}

var NetworkEvent []string = []string{
	"create",
	"connect",
	"disconnect",
	"destroy",
}

var VolumeEvent []string = []string{
	"create",
	"mount",
	"unmount",
	"destroy",
}

var DaemonEvent []string = []string{
	"reload",
	"shutdown",
}

var EvtAttribytes []string = []string{
	"name",                       //	컨테이너 이름
	"image",                      //	이미지 이름
	"exitCode",                   //	종료 코드
	"execDuration",               //	실행 시간
	"signal",                     //	kill signal
	"container",                  //	container ID
	"com.docker.compose.project", //	compose 프로젝트
	"com.docker.compose.service", //	compose 서비스
}

var EvtActionMap map[string][]string

func InitEventAction() {
	EvtActionMap = make(map[string][]string)

	for _, tp := range EventTypes {
		EvtActionMap[tp] = []string{}
		var events []string

		if tp == "container" {
			events = make([]string, len(ContainerEvent))
			copy(events, ContainerEvent)
		} else if tp == "image" {
			events = make([]string, len(ImagerEvent))
			copy(events, ImagerEvent)
		} else if tp == "network" {
			events = make([]string, len(NetworkEvent))
			copy(events, NetworkEvent)
		} else if tp == "volume" {
			events = make([]string, len(VolumeEvent))
			copy(events, VolumeEvent)
		} else if tp == "daemon" {
			events = make([]string, len(DaemonEvent))
			copy(events, DaemonEvent)
		}

		EvtActionMap[tp] = events
	}
}

func Contains(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return false
}

// FilterEvent는 이벤트 Type과 Action이 허용 목록에 있는지 확인
// 허용되면 true, 필터링(제외)되면 false 반환
func FilterEvent(evtType, evtAction string) bool {
	// EvtActionMap 초기화 확인
	if EvtActionMap == nil {
		InitEventAction()
	}

	// Type이 허용 목록에 있는지 확인
	actions, ok := EvtActionMap[evtType]
	if !ok {
		return false
	}

	// Action이 허용 목록에 있는지 확인
	return Contains(actions, evtAction)
}

// FilterAttrs는 허용된 Attribute만 추출
func FilterAttrs(attrs map[string]string) map[string]string {
	filtered := make(map[string]string)
	for _, key := range EvtAttribytes {
		if v, ok := attrs[key]; ok {
			filtered[key] = v
		}
	}
	return filtered
}
