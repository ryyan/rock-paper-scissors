import React from 'react';
import Websocket from 'react-websocket';

import Game from './Game';
import Scoreboard from './Scoreboard';
import CurrentPlayers from './Players';

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
    this.setState(JSON.parse(data));
  }

  render() {
    return (
      <div>
        <div className='pure-menu pure-menu-horizontal'>
          <h1 className='pure-menu-heading'>Rock Paper Scissors</h1>
          <CurrentPlayers />
        </div>

        <Game leftTaken={this.state.LeftTaken} rightTaken={this.state.RightTaken} />
        <Scoreboard wins={this.state.Wins} ties={this.state.Ties} games={this.state.PreviousGames} players={this.state.CurrentPlayers} />
        <Websocket url={'__WEBSOCKET_URL__/ws/game'} onMessage={this.handleData.bind(this)} />
      </div>
    );
  }
}
