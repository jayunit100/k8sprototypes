# PowerShell

To hack around on nodes with linux like fluency, you can use powershell.

- `ip a | grep 192` can be replaced with `Get-NetIPAddress | Select-String "192"`

# Windows networking

Things like hostNetwork and hostPort aren't really options.  If you need to access services such
as ingress controllers, you do so by: 

1) Implementing a NodePort service
2) DaemonSet for pods
3) Setting the NodePort Service to `externalTrafficPolicy=local`

	
- 
