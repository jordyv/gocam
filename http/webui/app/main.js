var data = {
    currentDate: new Date(),
    alerts: {}
}

function sortAlerts(alerts) {
    return alerts.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
}

fetch('/api/alerts')
    .then(r => r.json())
    .then(d => data.alerts = sortAlerts(d))

var app = new Vue({
    el: '#app',
    data: data,
    methods: {
        formatDate: function (dateString) {
            var date = new Date(dateString)
            var day = date.getDate(),
                month = date.getMonth() + 1,
                year = date.getFullYear(),
                hours = date.getHours(),
                minutes = date.getMinutes(),
                seconds = date.getSeconds()
            day = day < 10 ? `0${day}`: day
            month = month < 10 ? `0${month}`: month
            hours = hours < 10 ? `0${hours}`: hours
            minutes = minutes < 10 ? `0${minutes}`: minutes
            seconds = seconds < 10 ? `0${seconds}`: seconds
            return `${day}-${month}-${year} ${hours}:${minutes}:${seconds}`
        }
    }
})