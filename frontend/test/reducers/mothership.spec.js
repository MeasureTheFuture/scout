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
import expect from 'expect';
import { Mothership } from '../../src/reducers';

var locations = [{ "id":1,
                   "uuid":"800fd548-2d2b-4185-885d-6323ccbe88a0",
                   "ip_address":"192.168.0.1",
                   "authorised":true,
                   "name":"Chattanooga",
                   "State":"idle"},
                 { "id":2,
                   "uuid":"59ef7180-f6b2-4129-99bf-970eb4312b4b",
                   "ip_address":"192.168.0.1",
                   "authorised":true,
                   "name":"Brisbane",
                   "State":"calibrated"}]

describe('reducers', () => {
  describe('mothership', () => {
    it('should provide the initial state', () => {
      expect(
        Mothership(undefined, {})
      ).toEqual({locations:[], active:0})
    })

    it('should handle update locations', () => {
      expect(
        Mothership({locations:[], active:0}, {
          type:'UPDATE_LOCATIONS',
          locations:locations})
      ).toEqual({
        locations:locations,
        active:0
      })
    })

    it('should not allow a negative active index', () => {
      expect(
        Mothership({locations:locations, active:0}, {
          type:'SET_ACTIVE',
          active: -1
        })
      ).toEqual({
        locations:locations,
        active: 0
      })
    })

    it('should not allow an index larger than length of locations', () => {
      expect(
        Mothership({locations:locations, active: 0}, {
          type:'SET_ACTIVE',
          active: 3
        })
      ).toEqual({
        locations:locations,
        active: 1
      })
    })

    it('should set a valid index between 0 and the upper range', () => {
      expect(
        Mothership({locations:locations, active:0}, {
          type:'SET_ACTIVE',
          active: 1
        })
      ).toEqual({
        locations:locations,
        active: 1
      })
    })
  })
 })