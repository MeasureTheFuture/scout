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
import { UpdateActiveLocation, ActiveLocation } from '../reducers/index.js'
import PrimaryActions from './primaryActions.jsx';

var LocationEdit = React.createClass({
  render: function() {
    const { store } = this.context;

    return (
      <header className="locationLabel">
        <form className="pure-form">
          <h2 className="location-title"><input id="locationInput" className="location-title" type="text" placeholder={ActiveLocation(store).name} /></h2>
          <PrimaryActions />
        </form>
      </header>
    )
  }
});
LocationEdit.contextTypes = {
  store: React.PropTypes.object
}


var LocationLabel = React.createClass({
  render: function() {
    const { store } = this.context;

    return (
      <header className="locationLabel">
          <h2 className="location-title">{ActiveLocation(store).name}</h2>
          <PrimaryActions />
      </header>
    )
  }
});
LocationLabel.contextTypes = {
  store: React.PropTypes.object
}


var Placeholder = React.createClass({
  getFrameURL: function() {
    const { store } = this.context;

    if (!ActiveLocation(store).authorised) {
      return 'img/off-frame.gif';
    }

    if (ActiveLocation(store).state == 'calibrated') {
      return 'scouts/'+ActiveLocation(store).id+'/frame.jpg?d=' + new Date().getTime();
    } else if (ActiveLocation(store).state == 'calibrating') {
      return 'img/calibrating-frame.gif';
    }

    return 'img/uncalibrated-frame.gif';
  },

  render: function() {
    return (
      <div id="placeholder">
      <h3>&nbsp;</h3>
      <img className="pure-img placeholder" alt='test' src={this.getFrameURL()}/>
      </div>
    )
  }
});
Placeholder.contextTypes = {
  store: React.PropTypes.object
}

var Heatmap = React.createClass({
  toI: function(v) {
    return v | 0;
  },

  lerp: function(l, r, t) {
    return l + (r - l) * t
  },

  generateFill: function(t) {
    if (t < 0.001) {
       return "rgba(19, 27, 66, 0.1)"
    }

    if (t < 0.5) {
      return "rgba("+this.toI(this.lerp(19, 250, t))+","
        +this.toI(this.lerp(27, 212, t))+","
        +this.toI(this.lerp(66, 12, t))+","
        +this.lerp(0.3, 0.5, t)+")"
    }

    return "rgba("+this.toI(this.lerp(250, 186, t))+","
      +this.toI(this.lerp(212, 8, t))+","
      +this.toI(this.lerp(12, 16, t))+","
      +this.lerp(0.5, 0.5, t)+")"
  },

  maxTime: function(buckets) {
    var maxT = 0.0;

    buckets.map(function(i) {
      i.map(function(j) {
        maxT = Math.max(maxT, j);
      })
    })

    if (maxT < 0.1) {
      return 10.0;
    }

    return maxT;
  },

  pluralize :function(n) {
        return (n > 1) ? 's' : '';
  },

  secondsToStr: function(s) {
    s = Math.floor(s);
    var years = Math.floor(s / 31536000);
    if (years) {
        return years + ' yr' + this.pluralize(years);
    }

    var days = Math.floor((s %= 31536000) / 86400);
    if (days) {
        return days + ' day' + this.pluralize(days);
    }
    var hours = Math.floor((s %= 86400) / 3600);
    if (hours) {
        return hours + ' hr' + this.pluralize(hours);
    }
    var minutes = Math.floor((s %= 3600) / 60);
    if (minutes) {
        return minutes + ' min' + this.pluralize(minutes);
    }
    var seconds = s % 60;
    if (seconds) {
        return seconds + ' sec' + this.pluralize(seconds);
    }
    return '0 sec';
  },

  render: function() {
    const { store } = this.context;
    var url = 'scouts/'+ActiveLocation(store).id+'/frame.jpg?d=' + new Date().getTime();
    var buckets = ActiveLocation(store).summary.VisitTimeBuckets;
    var w = 1920;
    var h = 1080;
    var iBuckets = buckets.length;
    var jBuckets = buckets[0].length;
    var bucketW = w/iBuckets;
    var bucketH = h/jBuckets;
    var maxT = this.maxTime(buckets);
    var viewBox="0 0 " + w + " " + 1165;

    var data = []
    for (var i = 0; i < iBuckets; i++) {
      for (var j = 0; j < jBuckets; j++) {
        var t = 0.0;
        if (maxT > 0.0) {
          t = buckets[i][j] / maxT;
        }

        data.push(<rect key={i*iBuckets+j} x={i*bucketW} y={j*bucketH} width={bucketW} height={bucketH} style={{fill:this.generateFill(t)}} />);
      }
    }

    var scale = [];
    for (var i = 0; i < iBuckets; i++) {
      scale.push(<rect key={'s'+(i*iBuckets)} x={i*bucketW} y='1080' width={bucketW} height='30' style={{fill:this.generateFill((i*bucketW)/w)}} />);
    }

    return (
      <div id="heatmap">
      <h3>ACCUMULATED INTERACTION TIME</h3>
      <svg xmlns="http://www.w3.org/2000/svg" xmlnsXlink="http://www.w3.org/1999/xlink" viewBox={viewBox}>
        <image x="0" y="0" width={w} height={h} xlinkHref={url}/>
        {data}{scale}
        <text x="2" y={h+65} style={{fontSize:36,fontFamily:"Verdana, Geneva, sans-serif",fontWeight:"bold",letterSpacing:"-2px"}}>{this.secondsToStr(0.0)}</text>
        <text x="957" y={h+65} textAnchor="middle" style={{fontSize:36,fontFamily:"Verdana, Geneva, sans-serif",fontWeight:"bold",letterSpacing:"-2px"}}>{this.secondsToStr(0.5*maxT)}</text>
        <text x="1918" y={h+65} textAnchor="end" style={{fontSize:36,fontFamily:"Verdana, Geneva, sans-serif",fontWeight:"bold",letterSpacing:"-2px"}}>{this.secondsToStr(maxT)}</text>
      </svg>
      </div>
    )
  }
});
Heatmap.contextTypes = {
  store: React.PropTypes.object
}

var Analysis = React.createClass({
  render: function() {
    const { store } = this.context;
    var count = ActiveLocation(store).summary.VisitorCount;
    var vUpper = Math.ceil((count + 1) / 10) * 10;
    var vLower = vUpper - 10;


    var report = <p>No detected interactions.</p>;
    if (count > 0) {
        // have stayed for about <b>X</b> to <b>Y</b> minutes.</p>
        report = <p>Around <b>{vLower}</b> to <b>{vUpper}</b> visitors.</p>
    }

    return (
      <div id="analysis">
      <h3>INTERACTION REPORT</h3>
      {report}
      </div>
    )
  }
})
Analysis.contextTypes = {
  store: React.PropTypes.object
}

Location = React.createClass({
  render: function() {
    const { store } = this.context;
    var locationName = (store.getState().editLocation ? <LocationEdit /> : <LocationLabel /> )
    var heatmap = (ActiveLocation(store).state == 'measuring') ? <Heatmap /> : <Placeholder />
    var report = (ActiveLocation(store).state == 'measuring') ? <Analysis /> : " "

    return (
      <div className="location">
        <div id="locationName">{ locationName }</div>
        <div id="location-details">
          { heatmap }{ report }
        </div>
      </div>
    );
  }
});
Location.contextTypes = {
  store: React.PropTypes.object
};

export default Location;