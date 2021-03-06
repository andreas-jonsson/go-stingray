/*
Copyright (C) 2016-2017 Andreas T Jonsson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

/*
#include "sdk/engine_plugin_api/plugin_api.h"

typedef struct PluginApi PluginApi;
*/
import "C"

import (
	"unsafe"
)

type PluginApi C.PluginApi

var pluginApi PluginApi

var pluginName = C.CString("golang_plugin_loader")

//export GetName
func GetName() *C.char {
	return pluginName
}

//export get_plugin_api
func get_plugin_api(api C.unsigned) unsafe.Pointer {
	if api == C.PLUGIN_API_ID {
		f := GetName
		pluginApi.get_name = (*[0]byte)(unsafe.Pointer(&f))
		return unsafe.Pointer(&pluginApi)
	}
	return nil
}

func main() {
	// Not used
}
