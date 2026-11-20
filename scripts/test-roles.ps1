# Test raid-role RSVPs on the Guild Board (v0.8.33)
$ErrorActionPreference = "Stop"
$base = "http://localhost:4711/api/v1"
$tok = (Invoke-RestMethod -Uri "$base/tokens" -Method Post `
    -Headers @{ "X-Admin-Token" = "testadmin" } -ContentType "application/json" `
    -Body '{"label":"role-test"}').token
$h = @{ "Authorization" = "Bearer $tok" }

# create a test event
$ev = Invoke-RestMethod -Uri "$base/board/events" -Method Post -Headers $h `
    -ContentType "application/json" `
    -Body '{"title":"Test raid night","when_text":"Saturday 20:00","note":"role test"}'
$id = $ev.id
Write-Host "event id: $id"

# two members sign up with roles
Invoke-RestMethod -Uri "$base/board/events/$id/rsvp" -Method Post -Headers $h `
    -ContentType "application/json" -Body '{"name":"Brix","status":"yes","role":"heal"}' | Out-Null
Invoke-RestMethod -Uri "$base/board/events/$id/rsvp" -Method Post -Headers $h `
    -ContentType "application/json" -Body '{"name":"Emo","status":"yes","role":"dps"}' | Out-Null

# read back
$b = Invoke-RestMethod "$base/board"
$e = $b.events | Where-Object { $_.id -eq $id }
$rsvps = ($e.rsvps | ForEach-Object { "$($_.name)=$($_.status)/$($_.role)" }) -join ", "
Write-Host "rsvps: $rsvps"

# invalid role must be rejected
try {
    Invoke-RestMethod -Uri "$base/board/events/$id/rsvp" -Method Post -Headers $h `
        -ContentType "application/json" -Body '{"name":"Skumle","status":"yes","role":"healer"}'
    Write-Host "FAIL: bad role accepted"
} catch {
    $sr = New-Object IO.StreamReader($_.Exception.Response.GetResponseStream())
    Write-Host "bad role rejected: $($sr.ReadToEnd())"
}

# role change (heal -> tank) keeps status
Invoke-RestMethod -Uri "$base/board/events/$id/rsvp" -Method Post -Headers $h `
    -ContentType "application/json" -Body '{"name":"Brix","status":"yes","role":"tank"}' | Out-Null
$b2 = Invoke-RestMethod "$base/board"
$e2 = $b2.events | Where-Object { $_.id -eq $id }
$after = ($e2.rsvps | ForEach-Object { "$($_.name)=$($_.status)/$($_.role)" }) -join ", "
Write-Host "after change: $after"

# cleanup
Invoke-RestMethod -Uri "$base/board/events/$id" -Method Delete -Headers $h | Out-Null
Write-Host "test event cleaned up"