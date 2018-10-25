var data = {
    currentDate: new Date(),
    alertGroups: {},
    totalAlerts: 0,
}

function sortAlertGroups(groups) {
    var sortedGroupKeys = Object.keys(groups).sort((a, b) => new Date(b).getTime() - new Date(a).getTime())
    var sortedGroups = {}
    sortedGroupKeys.forEach(g => sortedGroups[g] = sortAlerts(groups[g]))
    return sortedGroups
}

function sortAlerts(alerts) {
    return alerts.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
}

function getItemCount(groups) {
    var count = 0
    Object.keys(groups).forEach(g => count += groups[g].length)
    return count
}

fetch('/api/alerts/grouped')
    .then(r => r.json())
    .then(d => {
        data.alertGroups = sortAlertGroups(d)
        data.totalAlerts = getItemCount(d)
    })

var app = new Vue({
    el: '#app',
    data: data,
    methods: {
        formatDateTime: function (dateString) {
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
        },
        formatDate: function (dateString) {
            var date = new Date(dateString)
            var day = date.getDate(),
                month = date.getMonth() + 1,
                year = date.getFullYear()
            day = day < 10 ? `0${day}`: day
            month = month < 10 ? `0${month}`: month
            return `${day}-${month}-${year}`
        },
        getTimestamp: function (dateString) {
            return (new Date(dateString)).getTime()
        },
        isFirstGroup: function (groupName) {
            return groupName === Object.keys(data.alertGroups)[0]
        }
    }
})