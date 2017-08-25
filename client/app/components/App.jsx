import React from 'react';
import Websocket from 'react-websocket';

import Game from './Game';
import Scoreboard from './Scoreboard';

export default class App extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      LeftTaken: false,
      RightTaken: false,
      Wins: [0, 0, 0],
      Ties: [0, 0, 0],
      PreviousGames: [{}]
    }
  }

  handleData(data) {
    let result = JSON.parse(data);
    this.setState(result);
    console.log(result);
  }

  render() {
    return (
      <div>
        <div className='pure-menu pure-menu-horizontal'>
          <h1 className='pure-menu-heading' href>Rock Paper Scissors</h1>
        </div>

        <Game leftTaken={this.state.LeftTaken} rightTaken={this.state.RightTaken} />
        <Scoreboard wins={this.state.Wins} ties={this.state.Ties} games={this.state.PreviousGames} />
        <Websocket url='ws://192.168.0.111:5000/websocket/rps'
          onMessage={this.handleData.bind(this)} />
      </div>
    );
  }
}

