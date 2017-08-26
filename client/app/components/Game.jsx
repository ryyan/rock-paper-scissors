import React from 'react';
import Websocket from 'react-websocket';

export default class Game extends React.Component {

  constructor(props) {
    super(props);
  }

  render() {
    return (
      <div id='game' className='pure-g'>
        <div id='left' className='pure-u-1-2'>
          <h1>{this.props.leftTaken ? '?' : '_'}</h1>
          <Button lor={'l'} choice={1} disabled={this.props.leftTaken} text='Rock' />
          <Button lor={'l'} choice={10} disabled={this.props.leftTaken} text='Paper' />
          <Button lor={'l'} choice={100} disabled={this.props.leftTaken} text='Scissors' />
        </div>
        <div id='right' className='pure-u-1-2'>
          <h1>{this.props.rightTaken ? '?' : '_'}</h1>
          <Button lor={'r'} choice={1} disabled={this.props.rightTaken} text='Rock' />
          <Button lor={'r'} choice={10} disabled={this.props.rightTaken} text='Paper' />
          <Button lor={'r'} choice={100} disabled={this.props.rightTaken} text='Scissors' />
        </div>
      </div>
    );
  }
}

class Button extends React.Component {

  constructor(props) {
    super(props);
    this.handleClick = this.handleClick.bind(this);
  }

  handleClick() {
    var request = new XMLHttpRequest();
    request.open('GET', 'http://192.168.0.111:5000/rps?lor=' + this.props.lor + '&choice=' + this.props.choice);
    request.send();
  }

  render() {
    let classname = 'pure-button pure-button-primary';
    if (this.props.disabled) {
      classname += ' pure-button-disabled';
    }

    return (
      <button className={classname} onClick={this.handleClick}>{this.props.text}</button>
    );
  }
}
