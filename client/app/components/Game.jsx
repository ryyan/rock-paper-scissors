import React from 'react';

export default class Game extends React.Component {

  render() {
    return (
      <div id='game' className='pure-g'>
        <div id='left' className='pure-u-1-2'>
          <h1>{this.props.leftTaken ? '???' : '___'}</h1>
          <Button leftOrRight={'l'} choice={1} disabled={this.props.leftTaken} text='Rock' />
          <Button leftOrRight={'l'} choice={10} disabled={this.props.leftTaken} text='Paper' />
          <Button leftOrRight={'l'} choice={100} disabled={this.props.leftTaken} text='Scissors' />
        </div>
        <div id='right' className='pure-u-1-2'>
          <h1>{this.props.rightTaken ? '???' : '___'}</h1>
          <Button leftOrRight={'r'} choice={1} disabled={this.props.rightTaken} text='Rock' />
          <Button leftOrRight={'r'} choice={10} disabled={this.props.rightTaken} text='Paper' />
          <Button leftOrRight={'r'} choice={100} disabled={this.props.rightTaken} text='Scissors' />
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
    let request = new XMLHttpRequest();
    request.open('POST', `__API_URL__/play?lor=${this.props.leftOrRight}&choice=${this.props.choice}`);
    request.send();
  }

  render() {
    let buttonClass = 'pure-button pure-button-primary';
    if (this.props.disabled) buttonClass += ' pure-button-disabled';

    return (
      <button className={buttonClass} onClick={this.handleClick}>{this.props.text}</button>
    );
  }
}
