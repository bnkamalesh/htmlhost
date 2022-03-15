(function () {
  const formDOM = document.getElementById("pagesubmit");
  const validationErrDOM = document.getElementById("validationErr");
  const a = parseInt(Math.random() * 100) % 10,
    b = parseInt(Math.random() * 100) % 10;
  const equationDOM = document.getElementById("equation");
  equationDOM.innerHTML = a + " + " + b;

  const mathDOM = document.getElementById("mathinput");
  mathDOM.onchange = () => {
    validationErrDOM.innerHTML = "";
    if (parseInt(mathDOM.value) == a + b) {
      mathDOM.className = "greenborder";
    } else {
      mathDOM.className = "redborder";
    }
  };

  const textareaDOM = formDOM.querySelector("textarea");
  textareaDOM.onchange = () => {
    validationErrDOM.innerHTML = "";
    if (!textareaDOM.value.trim()) {
      textareaDOM.className = "redborder";
    } else {
      textareaDOM.className = "greenborder";
    }
  };

  formDOM.onsubmit = (e) => {
    if (parseInt(mathDOM.value) !== a + b) {
      mathDOM.className = "redborder";
      validationErrDOM.innerHTML =
        'here you go, the correct answer is <strong style="font-style: normal">' +
        (a + b) +
        "<strong>";
      return false;
    }

    if (!textareaDOM.value.trim()) {
      textareaDOM.className = "redborder";
      validationErrDOM.innerHTML = "how about we add some content as well? :)";
      return false;
    }

    formDOM.setAttribute("action", "/#inputbody");
    formDOM.setAttribute("method", "POST");
  };

  // Dark/light theme selection
  const bodyDOM = document.getElementsByTagName("body")[0];
  let initialDark = false;
  if (
    window.matchMedia &&
    window.matchMedia("(prefers-color-scheme: dark)").matches
  ) {
    bodyDOM.className = "dark";
    initialDark = true;
  } else {
    bodyDOM.className = "";
  }

  window
    .matchMedia("(prefers-color-scheme: dark)")
    .addEventListener("change", (e) => {
      if (e.matches) {
        bodyDOM.className = "dark";
        initialDark = true;
      } else {
        bodyDOM.className = "";
      }
    });

  function getCookie(cookieName) {
    var name = cookieName + "=";
    var allCookieArray = document.cookie.split(";");
    for (var i = 0; i < allCookieArray.length; i++) {
      var temp = allCookieArray[i].trim();
      if (temp.indexOf(name) == 0)
        return temp.substring(name.length, temp.length);
    }
    return "";
  }

  switch (getCookie("userthemechoice")) {
    case "dark": {
      initialDark = true;
      bodyDOM.className = "dark";
      break;
    }
    case "light": {
      bodyDOM.className = "";
      initialDark = false;
      break;
    }
  }

  const darkToggleDOM = document.getElementById("darktoggle");
  darkToggleDOM.checked = initialDark;
  darkToggleDOM.addEventListener("change", (e) => {
    if (darkToggleDOM.checked) {
      bodyDOM.className = "dark";
      document.cookie = "userthemechoice=dark";
    } else {
      bodyDOM.className = "";
      document.cookie = "userthemechoice=light";
    }
  });

  const inputfiledom = document.getElementById("inputfile");
  inputfiledom.addEventListener("change", (e) => {
    if (!inputfiledom.files || !inputfiledom.files.length) {
      return;
    }

    const reader = new FileReader();
    reader.readAsText(inputfiledom.files[0], "UTF-8");
    reader.onload = function (e) {
      textareaDOM.value = e.target.result;
    };
  });

  const sse = (url, config = {}) => {
    const {
      onMessage,
      onError,
      initialBackoff = 1000, // milliseconds
      maxBackoff = 60 * 1000, // 60 seconds
      backoffStep = 1000, // milliseconds
    } = config;

    let backoff = initialBackoff,
      sseRetryTimeout = null;

    const start = () => {
      const source = new EventSource(url);
      const configState = { initialBackoff, maxBackoff, backoffStep, backoff };

      source.onopen = () => {
        // reset backoff to initial, so further failures will again start with initial backoff
        // instead of previous duration
        backoff = initialBackoff;
        configState.backoff = backoff;
      };

      source.onmessage = (event, configState) => {
        onMessage && onMessage(event, configState);
      };

      source.onerror = (err) => {
        source.close();
        clearTimeout(sseRetryTimeout);
        // reattempt connecting with *linear* backoff
        sseRetryTimeout = window.setTimeout(() => {
          start(url, onMessage);
          if (backoff < maxBackoff) {
            backoff += backoffStep;
            if (backoff > maxBackoff) {
              backoff = maxBackoff;
            }
          }
        }, backoff);
        onError && onError(err, configState);
      };
    };
    return start;
  };

  const clientID = Math.random()
    .toString(36)
    .replace(/[^a-z]+/g, "")
    .substring(0, 16);
  const activeClientsDOM = document.getElementById("live-viewers");
  const activePagesDOM = document.getElementById("live-pages");

  sse(`/sse/${clientID}`, {
    onMessage: (event) => {
      const { activeClients, activePages } = JSON.parse(event.data);
      activeClientsDOM.innerText = activeClients;
      activePagesDOM.innerText = activePages;
    },
    onError: (err) => {
      activeClientsDOM.innerText = "";
      activePagesDOM.innerText = "";
      console.log(err);
    },
  })();
})();

// drop/input an HTML file
function dropHandler(ev) {
  ev.preventDefault();
  console.log("File(s) dropped");

  if (ev.dataTransfer.items) {
    // Use DataTransferItemList interface to access the file(s)
    for (var i = 0; i < ev.dataTransfer.items.length; i++) {
      // If dropped items aren't files, reject them
      if (ev.dataTransfer.items[i].kind === "file") {
        var file = ev.dataTransfer.items[i].getAsFile();
        console.log("... file[" + i + "].name = " + file.name);
      }
    }
  } else {
    // Use DataTransfer interface to access the file(s)
    for (var i = 0; i < ev.dataTransfer.files.length; i++) {
      console.log(
        "... file[" + i + "].name = " + ev.dataTransfer.files[i].name
      );
    }
  }
}
