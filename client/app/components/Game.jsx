import React from 'react';
import Websocket from 'react-websocket';

export default class Game extends React.Component {

  constructor(props) {
    super(props);
  }

  render() {
    return (
      <div id='game' className='pure-u-1'>
        <div id='left' className='pure-u-1-3'>
          <h1>{this.props.leftTaken ? '?' : '_'}</h1>
        </div>
        <div id='mid' className='pure-u-1-3'>
          <h1>|</h1>
        </div>
        <div id='right' className='pure-u-1-3'>
          <h1>{this.props.rightTaken ? '?' : '_'}</h1>
        </div>
      </div>
    );
  }
}

