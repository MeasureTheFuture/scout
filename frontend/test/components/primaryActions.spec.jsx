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

import React from 'react';
import expect from 'expect';
import configureStore from 'redux-mock-store';
import { Provider } from 'react-redux';
import { mount } from 'enzyme';
import PrimaryActions from '../../src/components/PrimaryActions.jsx';
import jsdom from 'jsdom'

const doc = jsdom.jsdom('<!doctype html><html><body></body></html>')
global.document = doc
global.window = doc.defaultView
const mockStore = configureStore();

describe('components', () => {
  describe('PrimaryActions', () => {
	it('should provide an activate/deactivate button based on authorised state', () => {
		var s = {locations:[{id:2, uuid:'800fd548-2d2b-4185-885d-6323ccbe88a0', port:8080, authorised:true, name:'NYPL', state:'idle'}], active:0};

		var c = mount(<Provider store={mockStore(s)}><PrimaryActions /></Provider>);
		expect(c.find('#deactivate').length).toEqual(1);
		expect(c.find('#activate').length).toEqual(0);

		s.locations[s.active].authorised = false;
		c = mount(<Provider store={mockStore(s)}><PrimaryActions /></Provider>);
		expect(c.find('#deactivate').length).toEqual(0);
		expect(c.find('#activate').length).toEqual(1);
    })

    it('should only display a measure button if it is authorised and calibrated', () => {
		var s = {locations:[{id:2, uuid:'800fd548-2d2b-4185-885d-6323ccbe88a0', port:8080, authorised:true, name:'NYPL', state:'idle'}], active:0};
		var c = mount(<Provider store={mockStore(s)}><PrimaryActions /></Provider>);
		expect(c.find('#measure').length).toEqual(0);

		s.locations[s.active].authorised = false;
		s.locations[s.active].state = 'calibrated';
		c = mount(<Provider store={mockStore(s)}><PrimaryActions /></Provider>);
		expect(c.find('#measure').length).toEqual(0);

		s.locations[s.active].authorised = true;
		c = mount(<Provider store={mockStore(s)}><PrimaryActions /></Provider>);
		expect(c.find('#measure').length).toEqual(1);
    })

    it('should only display a calibrate button if it is authorised and idle or calibrated', () => {
		var s = {locations:[{id:2, uuid:'800fd548-2d2b-4185-885d-6323ccbe88a0', port:8080, authorised:false, name:'NYPL', state:'idle'}], active:0};
		var c = mount(<Provider store={mockStore(s)}><PrimaryActions /></Provider>);
		expect(c.find('#calibrate').length).toEqual(0);

		s.locations[s.active].state = 'calibrated';
		c = mount(<Provider store={mockStore(s)}><PrimaryActions /></Provider>);
		expect(c.find('#calibrate').length).toEqual(0);

		s.locations[s.active].authorised = true;
		c = mount(<Provider store={mockStore(s)}><PrimaryActions /></Provider>);
		expect(c.find('#calibrate').length).toEqual(1);

		s.locations[s.active].state = 'idle';
		c = mount(<Provider store={mockStore(s)}><PrimaryActions /></Provider>);
		expect(c.find('#calibrate').length).toEqual(1);
    })
  })
})
