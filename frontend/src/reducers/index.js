/*
 * Copyright (C) 2016 Clinton Freeman
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
"use strict;"

const initialState = {
  locations:[],
  active:0,
  editLocation:false,
  editSettings:false
}

function ActiveLocation(store) {
  var state = store.getState();
  return state.locations[state.active];
}

function GetLocations(store) {
    var httpreq = new XMLHttpRequest();
    httpreq.open("GET", "http://"+window.location.host+"/scouts", true);
    httpreq.send(null);
    httpreq.onreadystatechange = function() {
      if (httpreq.readyState == 4 && httpreq.status == 200) {
        var locations = JSON.parse(httpreq.responseText)
        store.dispatch({ type:'UPDATE_LOCATIONS', locations:locations})
      }
    }
}

function ClearMeasurements(store) {
  var state = store.getState();
  var l = Object.assign({}, state.locations[state.active]);

  var httpreq = new XMLHttpRequest();
  httpreq.open("GET", "http://"+window.location.host+"/scouts/"+l.uuid+"/clearMeasurements", true);
  httpreq.send(null);
  httpreq.onreadystatechange = function() {
    if (httpreq.readyState == 4 && httpreq.status == 200) {
      var locations = JSON.parse(httpreq.responseText)
      store.dispatch({ type:'UPDATE_LOCATIONS', locations:locations})
    }
  }
}

function SaveActiveLocation(store) {
  var state = store.getState();

  var l = Object.assign({}, state.locations[state.active]);

  // Push the active location to the backend.
  var httpreq = new XMLHttpRequest();
  httpreq.open("PUT", "http://"+window.location.host+"/scouts/"+l.uuid, true);
  httpreq.send(JSON.stringify(l));
  httpreq.onreadystatechange = function() {
    if (httpreq.readyState == 4 && httpreq.status == 200) {
      store.dispatch({ type:'UPDATE_LOCATIONS', locations:state.locations})
    }
  }
}

function UpdateActiveLocation(store, field, value) {
  var state = store.getState();

  var l = Object.assign({}, state.locations[state.active]);
  Reflect.set(l, field, value);
  state.locations[state.active] = l;

  SaveActiveLocation(store);
}

function Mothership(state, action) {
  if (state === undefined) {
    return initialState;
  }

  switch (action.type) {
    case 'UPDATE_LOCATIONS':
      return {
        locations: action.locations,
        active: state.active,
        editLocation: state.editLocation,
        editSettings: state.editSettings
      }

    case 'SET_ACTIVE':
      return {
        locations: state.locations,
        active: Math.min(state.locations.length - 1, Math.max(0, action.active)),
        editLocation: state.editLocation,
        editSettings: state.editSettings
      }

    case 'EDIT_LOCATION':
      return {
        locations: state.locations,
        active: state.active,
        editLocation: true,
        editSettings: state.editSettings
      }

    case 'SAVE_LOCATION':
      return {
        locations: state.locations,
        active: state.active,
        editLocation: false,
        editSettings: state.editSettings
      }

    case 'EDIT_SETTINGS':
      return {
        locations: state.locations,
        active: state.active,
        editLocation: state.editLocation,
        editSettings: true
      }

    case 'SAVE_SETTINGS':
      return {
        locations: state.locations,
        active: state.active,
        editLocation: state.editLocation,
        editSettings: false
      }

    default:
      return state;
  }
}

export { Mothership, GetLocations, ActiveLocation, UpdateActiveLocation, ClearMeasurements, SaveActiveLocation }
