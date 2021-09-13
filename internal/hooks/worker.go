package hooks

import (
	"fmt"
	"github.com/chen-keinan/go-user-plugins/uplugin"
	"github.com/chen-keinan/openshift-ordeal/pkg/models"
	"go.uber.org/zap"
	"plugin"
)

//PluginWorker instance which match command data to specific pattern
type PluginWorker struct {
	cmd *PluginWorkerData
	log *zap.Logger
}

//OpenshiftBenchAuditResultHook hold the plugin symbol for openshift bench audit result Hook
type OpenshiftBenchAuditResultHook struct {
	Plugins []plugin.Symbol
	Plug    *uplugin.PluginLoader
}

//NewPluginWorker return new plugin worker instance
func NewPluginWorker(commandMatchData *PluginWorkerData, log *zap.Logger) *PluginWorker {
	return &PluginWorker{cmd: commandMatchData, log: log}
}

//NewPluginWorkerData return new plugin worker instance
func NewPluginWorkerData(plChan chan models.OpenshiftAuditResults, hook OpenshiftBenchAuditResultHook, completedChan chan bool) *PluginWorkerData {
	return &PluginWorkerData{plChan: plChan, plugins: hook, completedChan: completedChan}
}

//PluginWorkerData encapsulate plugin worker properties
type PluginWorkerData struct {
	plChan        chan models.OpenshiftAuditResults
	completedChan chan bool
	plugins       OpenshiftBenchAuditResultHook
}

//Invoke invoke plugin accept audit bench results
func (pm *PluginWorker) Invoke() {
	go func() {
		ae := <-pm.cmd.plChan
		if len(pm.cmd.plugins.Plugins) > 0 {
			for _, pl := range pm.cmd.plugins.Plugins {
				_, err := pm.cmd.plugins.Plug.Invoke(pl, ae)
				if err != nil {
					pm.log.Error(fmt.Sprintf("failed to execute plugins %s", err.Error()))
				}
			}
		}
		pm.cmd.completedChan <- true
	}()
}
