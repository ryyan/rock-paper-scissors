import React from 'react';
import Websocket from 'react-websocket';

export default class Scoreboard extends React.Component {

  render() {
    return (
      <div id='scoreboard' className='pure-g'>
        <div id='scores' className='pure-u-1-2'>
          <Scores wins={this.props.wins} ties={this.props.ties} />
        </div>
        <div id='top10' className='pure-u-1-2'>
          <Games games={this.props.games} /> 
        </div>
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
          <tr>
            <td>Rock</td>
            <td>{this.props.wins[0]}</td>
            <td>{this.props.ties[0]}</td>
          </tr>
          <tr>
            <td>Paper</td>
            <td>{this.props.wins[1]}</td>
            <td>{this.props.ties[1]}</td>
          </tr>
          <tr>
            <td>Scissors</td>
            <td>{this.props.wins[2]}</td>
            <td>{this.props.ties[2]}</td>
          </tr>
        </tbody>
      </table>
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

