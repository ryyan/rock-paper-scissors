import React from 'react';
import Websocket from 'react-websocket';

export default class CurrentPlayers extends React.Component {


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
    return (
      <div>
        <p>Current Players: {this.state.players}</p>
        <Websocket url={'__WEBSOCKET_URL__/ws/players'} onMessage={this.handleData.bind(this)} />
      </div>
    );
  }
}
