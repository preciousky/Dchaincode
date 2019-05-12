var fs = require("fs");
var iconv = require('iconv-lite');


var fileStr = fs.readFileSync('initData.csv', { encoding: 'binary' });
var buf = new Buffer(fileStr, 'binary');
var str = iconv.decode(buf, 'GBK');

var resStr = "papers := []Paper{\n";

ConvertToTable(str, function (table, rowNum) {
    console.log(table);
    console.log("There are " + rowNum + " rows of instance in the table.")
    for (var i = 1; i <= rowNum; i++) {
        resStr += "Paper{";
        for (var j = 0; j < 26; j++) {
            resStr += table[0][j];
            resStr += ":\"";
            resStr += table[i][j];
            resStr += "\"";
            if (j < 25) {
                resStr += ",";
            }
        }
        resStr += "},\n";
    }
    resStr += "}";

    console.log(resStr);
})
fs.writeFileSync('./initCode.go', resStr);
console.log("程序执行完毕");

function ConvertToTable(data, callBack) {
    var table = new Array();
    var rows = new Array();
    rows = data.split("\r\n");
    var i = 0;
    for (i = 0; i < rows.length; i++) {
        table.push(rows[i].split(","));
    }
    callBack(table, i - 2);
}