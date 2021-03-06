/*
 * Copyright (C) 2017 Clinton Freeman
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
import React from 'react';
import { UpdateActiveLocation, ActiveLocation, GetLocations, SaveActiveLocation } from '../reducers/index.js';

var Settings = React.createClass({
  updateField: function(field) {
  	const { store } = this.context;
  	var state = store.getState();

  	var v = document.getElementById(field).value;
  	if (v) {
  		var l = Object.assign({}, state.locations[state.active]);
  		Reflect.set(l, field, Number(document.getElementById(field).value));
  		state.locations[state.active] = l;
  	}
  },

  handleSave: function() {
  	const { store } = this.context;
  	var state = store.getState();

  	this.updateField("MinArea");
    this.updateField("MaxArea");
  	this.updateField("MinDuration");
  	this.updateField("IdleDuration");
  	this.updateField("MogHistoryLength");
  	SaveActiveLocation(store);

    store.dispatch({ type:'SAVE_SETTINGS' })
  },

  handleCancel: function() {
  	const { store } = this.context;

  	console.log("cancel");
  	store.dispatch({ type:'SAVE_SETTINGS' })
  },

  render: function() {
  	const { store } = this.context;

    return (
      <form className="pure-form pure-form-stacked">
      <fieldset>
      <div className="pure-g">
        <div className="pure-u-5-5"><h3>SETTINGS</h3></div>
      </div>
        <div class="pure-control-group">
          <label htmlFor="MinArea">Minimum Area: &nbsp;</label>
          <input id="MinArea" type="text" placeholder={ActiveLocation(store).MinArea}></input>
          <span className="pure-form-message">The minimum area (in pixels) of a detected object before it gets counted as a person.</span>
        </div>

        <div class="pure-control-group">
          <label htmlFor="MaxArea">Maximum Area: &nbsp;</label>
          <input id="MaxArea" type="text" placeholder={ActiveLocation(store).MaxArea}></input>
          <span className="pure-form-message">The maximum area (in pixels) of a detected object before it gets counted as a person.</span>
        </div>

        <div class="pure-control-group">
          <label htmlFor="MinDuration">Minimum Duration: &nbsp;</label>
          <input id="MinDuration" type="text" placeholder={ActiveLocation(store).MinDuration}></input>
          <span className="pure-form-message">The minimum time (in seconds) a detected object must be tracked before it gets counted as a person.</span>
        </div>

        <div class="pure-control-group">
          <label htmlFor="IdleDuration">Idle Duration: &nbsp;</label>
          <input id="IdleDuration" type="text" placeholder={ActiveLocation(store).IdleDuration}></input>
          <span className="pure-form-message">If an object is briefly occluded, tracking can be resumed. IdleDuration is the maximum time (in seconds) that a detected object can be 'resumed'.</span>
        </div>

        <div class="pure-control-group">
          <label htmlFor="MogHistoryLength">History Length: &nbsp;</label>
          <input id="MogHistoryLength" type="text" placeholder={ActiveLocation(store).MogHistoryLength}></input>
          <span className="pure-form-message">The number of frames to be used when calculating the background frame for the subtraction algorithm.</span>
        </div>

        <div class="pure-controls">
            <a className="pure-button pure-button-primary" href="#" onClick={this.handleSave}>save</a>
            <span className="pure-form-message-inline"><a href="#" onClick={this.handleCancel}>cancel</a></span>
        </div>
      </fieldset>
      </form>
    )
  }
});
Settings.contextTypes = {
  store: React.PropTypes.object
}

export default Settings;