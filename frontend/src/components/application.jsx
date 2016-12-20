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
"use strict";

import React from 'react';
import { GetLocations } from '../reducers/index.js'
import Location from './location.jsx';
import Introduction from './documentation.jsx';

var NavItem = React.createClass({
  handleClick: function() {
    store.dispatch({ type:'SET_ACTIVE', active:this.props.idx})
  },

  render: function() {
    return (
      <li className="navItem">
        <a href="#" onClick={this.handleClick}>[{this.props.name}]</a>
      </li>
    );
  }
});

var NavList = React.createClass({
  render: function() {
    var navNodes = ""
    if (this.props.data.map.length > 1) {
      navNodes = this.props.data.map(function(location, index) {
        return (
          <NavItem name={location.name} key={index} idx={index} />
        )
      });
    }

    return (
      <ul className="navList">
        {navNodes}
        <li className="navItem">&nbsp;</li>
        <li className="navItem"><a href="/download.zip">[<i className="fa fa-download"></i> Download Data]</a></li>
      </ul>
    )
  }
})

var Application = React.createClass({
  loadFromServer: function () {
    const { store } = this.context;
    GetLocations(store);
  },

  componentDidMount: function() {
    this.loadFromServer();
    setInterval(this.loadFromServer, 1000);
  },

  render: function() {
    const { store } = this.context;

    var state = store.getState();
    var mainContent = ((state.locations.length) ? <Location /> : <Introduction />);

    return (
      <div className="pure-g">
        <div className="sidebar pure-u-1 pure-u-md-1-4">
          <div className="header">
            <h1 className="brand"><img className="pure-img" alt='Measure the Future logo' src='/img/logo.gif' /></h1>
            <nav className="nav"><NavList data={state.locations} /></nav>
          </div>
        </div>
        <div className="content pure-u-1 pure-u-md-3-4" id="results">
          {mainContent}
        </div>
      </div>
    )
  }
})
Application.contextTypes = {
  store: React.PropTypes.object
};

export default Application