package main

const html = `<!doctype html>
<html lang="en">
	<head>
	<!-- Required meta tags -->
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

	<!-- Bootstrap CSS -->
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
	<title>Employees App</title>

	</head>
	<body>
		<div class="container-fluid h-100">
			<div class="row bg-veeam-light align-items-center">
				<div class="col-lg-1"></div>
				<div class="col-lg-8   rounded border border-light">
					<div class="m-4">
					<h1>Employees</h1>
					</div>
				</div>
			</div>
			<div class="row" id="employees" >
					<div class="col-lg-1"></div>
                    <div class="col-lg-8">
                      <table class="table">
					    <tr>
							<th></th>
							<th>No</th>
							<th>Hire Date</th>
							<th>First Name</th>
							<th>Last Name</th>
							<th>Birth Date</th>
							<th>Gender</th>
						</tr>
						<tr>
							<td><button class="btn btn-success"  onclick="add();">Add</button></td>
							<td><input type="number" value=1 id="iemp_no"></input></td>
							<td><input type="text" id="ihire_date" value="2000-04-16"></input></td>
							<td><input type="text" id="ifirst_name"></input></td>
							<td><input type="text" id="ilast_name"></input></td>
							<td><input type="text" id="ibirth_date" value="2020-04-16"></input></td>
							<td><input type="text" id="igender" value="X"></input></td>
						</tr>
                        <tr v-for="emp in employees"   >
						  <td><button class="btn btn-danger rounded-circle" v-bind:id="emp.emp_no" onclick="deleterec(this.id);">X</button></td>
						  <td>{{emp.emp_no}}</td>
						  <td>{{emp.hire_date | formatGoDate}}</td>
                          <td>{{emp.first_name}}</td>
						  <td>{{emp.last_name}}</td>
						  <td>{{emp.birth_date | formatGoDate}}</td>
						  <td>{{emp.gender}}</td>
                        </tr>
                      </table>
                    </div>
            </div>
		</div>
				
		<!-- Optional JavaScript -->
		<!-- jQuery first, then Popper.js, then Bootstrap JS -->
		<script src="https://code.jquery.com/jquery-3.5.1.min.js" integrity="sha384-ZvpUoO/+PpLXR1lu4jmpXWu80pZlYUAfxl5NsBMWOEPSjUn/6Z/hRTt8+pR6L4N2" crossorigin="anonymous"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.12.9/umd/popper.min.js" integrity="sha384-ApNbgh9B+Y1QKtv3Rn7W3mgPxhU9K/ScQsAP7hUibX39j7fakFPskvXusvfa0b4Q" crossorigin="anonymous"></script>
		<script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/js/bootstrap.min.js" integrity="sha384-JZR6Spejh4U02d8jOt6vLEHfe/JQGiRRSQQxSfFWpi1MquVdAyjUar5+76PVCmYl" crossorigin="anonymous"></script>
		<script src="https://cdn.jsdelivr.net/npm/vue@2"></script>
		<script>
			window.resturl = location.href.split('?')[0].replace(/\\/$/,'')+'/api'

			Vue.filter('formatGoDate', function(value) {
				if (value) {
					return (new Date(Date.parse(value))).toDateString();
				} else {
					return ""
				}
			})

			window.employees = new Vue({
				el: "#employees",
				data: {
					employees : [
					]   
				}
			})

			//god i hate golang crappy time convers while demarshalling
			var offset = (new Date()).getTimezoneOffset()
			var offseth = parseInt(offset/60)
			var offsetabs = Math.abs(offseth)
			offsetabs = (offsetabs<10)?"0"+offsetabs:""+offsetabs
			var offsetm = Math.abs(offset%60)
			offsetm = (offsetm<10)?"0"+offsetm:""+offsetm
			window.tz = (((offset/60)<0)?"-":"+")+offsetabs+":"+offsetm
			window.dateappend = "T00:00:00"+window.tz

			function add() {
				var empnow = parseInt($("#iemp_no").val())
				$("#iemp_no").val(empnow+1)


				fetch(window.resturl, {
					method: 'POST', // *GET, POST, PUT, DELETE, etc.
					mode: 'cors', // no-cors, *cors, same-origin
					cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
					credentials: 'same-origin', // include, *same-origin, omit
					redirect: 'follow', // manual, *follow, error
					referrerPolicy: 'no-referrer', // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
					body: JSON.stringify({emp_no:empnow,
						hire_date:$("#ihire_date").val()+window.dateappend,
						birth_date:$("#ibirth_date").val()+window.dateappend,
						first_name:$("#ifirst_name").val(),
						last_name:$("#ilast_name").val(),
						gender:$("#igender").val()})
				}).then(response => response).then(data => {
					console.log("added")
				}).catch((error) => {
					console.error('Error:', error);
				}).finally(()=> {
				})
			}
			function deleterec(id) {
				fetch(window.resturl, {
					method: 'DELETE', // *GET, POST, PUT, DELETE, etc.
					mode: 'cors', // no-cors, *cors, same-origin
					cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
					credentials: 'same-origin', // include, *same-origin, omit
					redirect: 'follow', // manual, *follow, error
					referrerPolicy: 'no-referrer', // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
					body: JSON.stringify({deleteid:parseInt(id)})
				}).then(response => response).then(data => {
					console.log("deleted")
				}).catch((error) => {
					console.error('Error:', error);
				}).finally(()=> {
				})
			}

			function refresh() {
				console.log("Refreshing")
				fetch(window.resturl, {
					method: 'GET', // *GET, POST, PUT, DELETE, etc.
					mode: 'cors', // no-cors, *cors, same-origin
					cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
					credentials: 'same-origin', // include, *same-origin, omit
					redirect: 'follow', // manual, *follow, error
					referrerPolicy: 'no-referrer' // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
				}).then(response => response.json()).then(data => {
					window.employees.employees = []
					var upl = 1
					data.forEach(emp => {
						var empno = parseInt(emp.emp_no)
						if (empno > upl) {
							upl = empno
						}
						window.employees.employees.push(emp)
					})
					$("#iemp_no").val(upl+1)
					setTimeout(refresh, 1000);
				}).catch((error) => {
					console.error('Error:', error," delaying refresh 10s");
					setTimeout(refresh, 10000);
				}).finally(()=> {
					
				})
			}
			refresh();


		</script>
	</body>
</html>
`
