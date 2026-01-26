package docker

type Type string

var eventTypes []string = []string{
	"container",
	"daemon",
	"image",
	"network",
	"volume",
}

var containerEvent []string = []string{
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

var imagerEvent []string = []string{
	"pull",
	"push",
	"tag",
	"untag",
	"delete",
	"save",
	"load",
}

var networkEvent []string = []string{
	"create",
	"connect",
	"disconnect",
	"destroy",
}

var volumeEvent []string = []string{
	"create",
	"mount",
	"unmount",
	"destroy",
}

var daemonEvent []string = []string{
	"reload",
	"shutdown",
}

var evtAttribytes []string = []string{
	"name",                       //	컨테이너 이름
	"image",                      //	이미지 이름
	"exitCode",                   //	종료 코드
	"execDuration",               //	실행 시간
	"signal",                     //	kill signal
	"container",                  //	container ID
	"com.docker.compose.project", //	compose 프로젝트
	"com.docker.compose.service", //	compose 서비스
}

var evtActionMap map[string][]string

func initEventAction() {
	evtActionMap = make(map[string][]string)

	for _, tp := range eventTypes {
		evtActionMap[tp] = []string{}
		var events []string

		if tp == "container" {
			events = make([]string, len(containerEvent))
			copy(events, containerEvent)
		} else if tp == "image" {
			events = make([]string, len(imagerEvent))
			copy(events, imagerEvent)
		} else if tp == "network" {
			events = make([]string, len(networkEvent))
			copy(events, networkEvent)
		} else if tp == "volume" {
			events = make([]string, len(volumeEvent))
			copy(events, volumeEvent)
		} else if tp == "daemon" {
			events = make([]string, len(daemonEvent))
			copy(events, daemonEvent)
		}

		evtActionMap[tp] = events
	}
}

func contains(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return false
}
