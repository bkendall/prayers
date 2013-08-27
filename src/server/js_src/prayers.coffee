placePrayers = (html) ->
    $("div#prayers").html(html)

loadData =
    type: "POST"
    url: "//prayers.bryankendall.com/load"
    async: false
    data:
        Integration: $("div#prayers").data("integration")
    success: placePrayers
$.ajax(loadData)

PrayerSubmitHandler = (event) ->
    event.preventDefault()
    postData =
        type: "POST"
        url: "//prayers.bryankendall.com/submit"
        data:
            Integration: $("div#prayers").data("integration")
            FirstName: $("input[name='FirstName']").val()
            LastName: $("input[name='LastName']").val()
            Email: $("input[name='Email']").val()
            Title: $("input[name='Title']").val()
            Body: $("textarea[name='Body']").val()
    $.ajax(postData)

PrayHandler = (event) ->
    event.preventDefault()
    postData =
        type: "POST"
        url: "//prayers.bryankendall.com/pray"
        data:
            Integration: $("div#prayers").data("integration")
            Id: $(event.currentTarget).data("id")
    $.ajax(postData)

PrayerSubmit = $("a#SubmitPrayer")
PrayerSubmit.on("click", PrayerSubmitHandler)
$("a.pray").on("click", PrayHandler)
