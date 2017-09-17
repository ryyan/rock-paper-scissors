// See http://brunch.io for documentation.
exports.files = {
  javascripts: {
    joinTo: {
      'vendor.js': /^(?!app)/,
      'app.js': /^app/
    }
  },

  stylesheets: {joinTo: 'app.css'}
};

exports.npm = {
  enabled: true,
  styles: {
    purecss: ['build/pure-min.css']
  }
};

exports.plugins = {
  babel: {
    presets: ['latest', 'react']
  },

  replacer: {
    dict: [
      {
        key: '__API_URL__',
        value: 'http://localhost:5000'
      },
      {
        key: '__WEBSOCKET_URL__',
        value: 'ws://localhost:5001'
      }
    ]
  }
};

