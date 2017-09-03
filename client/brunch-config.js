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

exports.plugins = {
  replacer: {
    dict: [
      {
        key: '__API_HOST__',
        value: 'localhost:5000'
      }
    ]
  },

  babel: {
    presets: ['latest', 'react']
  }
};

exports.npm = {
  enabled: true,
  styles: {
    purecss: ['build/pure-min.css']
  }
};
