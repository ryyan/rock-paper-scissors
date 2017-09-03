import React from 'react';
import Websocket from 'react-websocket';

export default class Scoreboard extends React.Component {

  render() {
    return (
      <div id='scoreboard' className='pure-g'>
        <div id='scores' className='pure-u-1-2'>
          <Scores wins={this.props.wins} ties={this.props.ties} />
          <CurrentPlayers />
        </div>
        <div id='top10' className='pure-u-1-2'>
          <Games games={this.props.games} /> 
        </div>
      </div>
    );
  }
}

class CurrentPlayers extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      players: 0
    }
  }

  handleData(data) {
    this.setState({players: data});
  }

  render() {
    let currentPlayersWsUrl = 'ws://' + '__API_HOST__' + '/rps/ws/players';

    return (
      <div>
        <p>Current Players: {this.state.players}</p>
        <Websocket url={currentPlayersWsUrl} onMessage={this.handleData.bind(this)} />
      </div>
    );
  }
}

class Scores extends React.Component {

  render() {
    return (
      <table className='pure-table'>
        <thead>
          <tr>
            <th></th>
            <th>Wins</th>
            <th>Ties</th>
          </tr>
        </thead>
        <tbody>
          <Score name={'Rock'} wins={this.props.wins[0]} ties={this.props.ties[0]} />
          <Score name={'Paper'} wins={this.props.wins[1]} ties={this.props.ties[1]} />
          <Score name={'Scissors'} wins={this.props.wins[2]} ties={this.props.ties[2]} />
        </tbody>
      </table>
    );
  }
}

class Score extends React.Component {

  render () {
    return (
      <tr>
        <td>{this.props.name}</td>
        <td>{this.props.wins}</td>
        <td>{this.props.ties}</td>
      </tr>
    );
  }
}

class Games extends React.Component {

  render() {
    let gameRows = this.props.games.reverse().map((game) =>
      <tr>
        <td>{game.Left}</td>
        <td>{game.Right}</td>
      </tr>
    );

    return (
      <table className='pure-table'>
        <thead>
          <tr>
            <th></th>
            <th>Last 10</th>
          </tr>
        </thead>
        <tbody>
          {gameRows}
        </tbody>
      </table>
    );
  }
}

