(function() {
  function ready(fn) {
    if (document.readyState !== 'loading') {
      fn();
    } else {
      document.addEventListener('DOMContentLoaded', fn);
    }
  }

  function connectCalendar() {
    const calendar = document.getElementById('report-date-selector');
    calendar.addEventListener('change', (event) => {
      window.location.search = `?date=${event.target.value}`;
    })

    const sp = new URLSearchParams(window.location.search);
    const date = sp.get('date');
    if (date) {
      calendar.value = date;  
    }
  }

  ready(connectCalendar);
})();