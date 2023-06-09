document.getElementById("generate-report").addEventListener("click", function() {
    fetch("/report")
        .then(function(response) {
            return response.json();
        })
        .then(function(data) {
            var reportContainer = document.getElementById("report-container");
            reportContainer.innerHTML = "<h2>" + data.title + "</h2><p>" + data.content + "</p>";
        })
        .catch(function(error) {
            console.log("Error:", error);
        });
});
function getHosts() {
    fetch("/hosts")
        .then(function(response) {
            return response.json();
        })
        .then(function(data) {
            var tableBody = document.querySelector("#host-table tbody");
            tableBody.innerHTML = "";

            data.forEach(function(host) {
                var row = document.createElement("tr");

                var idCell = document.createElement("td");
                idCell.textContent = host.id;
                row.appendChild(idCell);

                var nameCell = document.createElement("td");
                nameCell.textContent = host.name;
                row.appendChild(nameCell);

                var ipCell = document.createElement("td");
                ipCell.textContent = host.ip;
                row.appendChild(ipCell);

                var locationCell = document.createElement("td");
                locationCell.textContent = host.location;
                row.appendChild(locationCell);

                var actionCell = document.createElement("td");
                var deleteButton = document.createElement("button");
                deleteButton.textContent = "Delete";
                deleteButton.addEventListener("click", function() {
                    deleteHost(host.id);
                });
                actionCell.appendChild(deleteButton);
                row.appendChild(actionCell);

                tableBody.appendChild(row);
            });
        })
        .catch(function(error) {
            console.log("Error:", error);
        });
}

function deleteHost(id) {
    fetch("/deletehost?id=" + id, { method: "DELETE" })
        .then(function(response) {
            if (response.ok) {
                getHosts();
            } else {
                console.log("Delete host failed.");
            }
        })
        .catch(function(error) {
            console.log("Error:", error);
        });
}

document.addEventListener("DOMContentLoaded", function() {
    getHosts();
});

// script.js

function scanHosts() {
    // 发起扫描主机请求
    fetch("/scanhosts")
        .then(response => {
            if (response.ok) {
                // 扫描成功
                return response.json();
            } else {
                // 扫描失败
                throw new Error("Scan Hosts Failed");
            }
        })
        .then(data => {
            // 构造扫描结果的HTML输出
            let output = "<h2>Scan Results:</h2>";
            output += "<table>";
            output += "<tr><th>ID</th><th>Name</th><th>IP</th><th>Status</th></tr>";
            data.forEach(host => {
                output += "<tr>";
                output += "<td>" + host.id + "</td>";
                output += "<td>" + host.name + "</td>";
                output += "<td>" + host.ip + "</td>";
                output += "<td>" + host.status + "</td>";
                output += "</tr>";
            });
            output += "</table>";

// 显示扫描结果
            document.getElementById("scanResults").innerHTML = output;
        })
        .catch(error => {
            console.log("Scan Hosts Error:", error);
        });
}
// JavaScript code for handling dynamic updates and form submission

// Get the form element
const addHostForm = document.getElementById('addHostForm');

// Get the host list container
const hostList = document.getElementById('hostList');

// Add an event listener for form submission
addHostForm.addEventListener('submit', function(event) {
    event.preventDefault(); // Prevent form submission

    // Get the input values
    const nameInput = document.getElementById('nameInput');
    const ipInput = document.getElementById('ipInput');
    const locationInput = document.getElementById('locationInput');

    // Create a new host element
    const hostElement = document.createElement('div');
    hostElement.classList.add('host');

    // Create the host details
    const nameElement = document.createElement('span');
    nameElement.textContent = 'Name: ' + nameInput.value;
    hostElement.appendChild(nameElement);

    const ipElement = document.createElement('span');
    ipElement.textContent = 'IP: ' + ipInput.value;
    hostElement.appendChild(ipElement);

    const locationElement = document.createElement('span');
    locationElement.textContent = 'Location: ' + locationInput.value;
    hostElement.appendChild(locationElement);

    // Add the host to the host list
    hostList.appendChild(hostElement);

    // Clear the input fields
    nameInput.value = '';
    ipInput.value = '';
    locationInput.value = '';
});




