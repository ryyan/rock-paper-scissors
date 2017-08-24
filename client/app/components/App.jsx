import React from 'react';
import Websocket from 'react-websocket';

export default class App extends React.Component {

  handleData(data) {
    let result = JSON.parse(data);
    console.log(data);
  }

  render() {
    return (
      <div>
        <Websocket url='ws://192.168.0.111:5000/websocket/rps'
          onMessage={this.handleData.bind(this)} />
      </div>
    );
  }
}
