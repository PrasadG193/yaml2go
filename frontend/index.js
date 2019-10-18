//const URL = "http://localhost:8080/v1/convert"
const URL = "https://us-central1-yaml2go.cloudfunctions.net/yaml2go"

let go = document.getElementById("goGenerator")

let editor = ""

window.generatorCall = function () {
  let yamlData = document.getElementById("yamlGenerator").value
  document.getElementById('yamlGenerator').style.border = "1px solid #ced4da"
  yamlData = editor.getValue()
  $.ajax({
    'url': `${URL}`,
    'type': 'POST',
    'data': yamlData,
    'success': function (data) {
      go.setValue(data)
    },
    'error': function (jqXHR, request, error) {
      document.getElementById('yamlGenerator').style.border = "1px solid red"
      if (jqXHR.status == 400) {
        alert('Invalid yaml format')
      } else {
        alert('Something went wrong! Please report this to me@prasadg.dev')
      }
    }
  });

}

//Convert
document.getElementById("convert").addEventListener('click', () => {
  console.log('Temp Data')
  generatorCall()
})

//Clear YAML
document.getElementById('clearYaml').addEventListener('click', () => {
  editor.setValue('')
})

//Clear Go
document.getElementById('clearGo').addEventListener('click', () => {
  go.setValue('')
})


$(document).ready(function () {
  //code here...
  var input = $(".codemirror-textarea")[0];
  var output = $(".codemirror-textarea")[1];
  editor = CodeMirror.fromTextArea(input, {
    mode: "text/x-yaml",
    lineNumbers: true
  });
  editor.setSize(600, 400)

  go = CodeMirror.fromTextArea(output, {
    lineNumbers: true,
    mode: "text/x-go"
  });
  go.setSize(600, 400)
  $('.CodeMirror').resizable({
    resize: function () {
      editor.setSize($(this).width(), $(this).height());
    }
  });
});
