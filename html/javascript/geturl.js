function downURL(output,url){
    var req = new window.XMLHttpRequest();
    req.onreadystatechange = function(){
      if(req.readyState == 4 && req.status == 200){
        var data=req.responseText;
        // data = JSON.parse(data);
        console.log(data);
        document.getElementById(output).innerHTML = data;
        // document.getElementById(output).innerHTML = data;
      }
    };
    req.open("POST","/download",true);
    req.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    var str = "url="+url
    req.send(str);
}

function getMp3List(output) {

  var req = new window.XMLHttpRequest();
  req.onreadystatechange = function(){
    if(req.readyState == 4 && req.status == 200){
      var data=req.responseText;
      data = JSON.parse(data);
      console.log(data)
      document.getElementById(output).innerHTML = createDownloadLink(data);
    }
  };
  req.open("GET","/list",true);
  req.send( null );
}

function createDownloadLink(data) {
  var output = ""
  for(var i=0;i<data.length;i++) {
    var tmp = data[i]
    if (i>0) {
      output += "<br>"
    }
    output += "<a href='/mp3/"+tmp.No+"'>"+tmp.Name+"<a>"
  }
  return output
}

var health_down_flag = false;
var health_down_url =[];
var health_down_cmd ="";

function healthck(output) {
  var req = new window.XMLHttpRequest();
  req.onreadystatechange = function(){
    if(req.readyState == 4 && req.status == 200){
      var data=req.responseText;
      data = JSON.parse(data);
      health_down_flag = data.Downflag
      health_down_url = data.DownUrl
      health_down_cmd = data.DownCmd
      if (health_down_flag) {
        document.getElementById(output).innerHTML = health_down_cmd;
      }else{
        document.getElementById(output).innerHTML = "not down load";
        
      }
      // document.getElementById(output).innerHTML = data;
    }
  };
  req.open("GET","/health",true);
  req.send( null );
}

var timer = null;

function watch_cmd(output) {
  if (timer != null)  {
    clearInterval(timer);
  }else{
    console.log( "timerstart");
  }
  healthck(output)
  if (health_down_flag) {
    if ((health_down_cmd != "")&&(health_down_cmd != null)){
      document.getElementById(output).innerHTML = health_down_cmd;
    }
    timer = setTimeout(watch_cmd,2000,output)
  }else {
    document.getElementById(output).innerHTML = "Check End";
    console.log("TimeEnd");
    timer = null;
  }
}