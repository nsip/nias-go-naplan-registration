    // <!-- set up the sse connections  -->
    var txID = "unknown";
    var feedSource = new EventSource("/validate/statusfeed/unknown");
    var readySource = new EventSource("/validate/readyfeed/unknown");




    window.onload = function() {

        $("#fetch").prop('disabled', true);
        $("#fetch2").prop('disabled', true);

        var fileInput = document.getElementById('fileInput');
        // var fileDisplayArea = document.getElementById('fileDisplayArea');

        fileInput.addEventListener('change', function(e) {
            $("#fetch").prop('disabled', true);
            $("#fetch2").prop('disabled', true);
            var file = fileInput.files[0];
            var textType = /text.*/;
            var textTypeMS = /application.*excel/;

			// Temporary remove mimetype - not supported on some Windows Browsers
            // if (file.type.match(textType) || file.type.match(textTypeMS)) {

				var fileName = file.name.replace( /[<>:"\/\\|?*]+/g, '' );
                var reader = new FileReader();
                var progressNode = document.getElementById("upload-progress");

                reader.onprogress = function(event) {
                    if (event.lengthComputable) {
                        progressNode.max = event.total;
                        progressNode.value = event.loaded;
                    }
                };

                reader.onload = function(e) {
                    // fileDisplayArea.innerText = reader.result;
                    $.post("/validate/naplan/reg/" + fileName, reader.result, function(data) {
                        txID = data
                            // console.log(txID)

                        ref = "/validate/statusfeed/" + txID;
                        feedSource.close();
                        feedSource = new EventSource(ref);

                        feedSource.onmessage = function(event) {
                            var data = JSON.parse(event.data);

                            $("#result").empty();

                            var svg = dimple.newSvg("#result", 400, 200);
                            var chart = new dimple.chart(svg, data);
                            chart.addCategoryAxis("y", "v_type");
                            chart.addMeasureAxis("x", "count");
                            chart.addSeries(null, dimple.plot.bar);
                            chart.draw();

                        };

                        ref = "/validate/readyfeed/" + txID
                        readySource.close();
                        readySource = new EventSource(ref);

                        readySource.onmessage = function(event) {

                            if (txID == event.data) {
                                $("#fetch").prop('disabled', false);
                                $("#fetch2").prop('disabled', false);
                            }

                        };


                        // document.getElementById("uploadresult").innerHTML = txID;
                        document.getElementById("uploadresult").innerHTML = "File Accepted";
                    });
                }

                // reader.readAsText(file);    
                reader.readAsBinaryString(file);
            // } else {
            //     document.getElementById("uploadresult").innerHTML = "File Type Not Supported";
            // }
        });
    }



    function renderAnalysis(txID) {
        // var errorsBarChart = dc.barChart("#errors-chart");
        // var errorsByTypeChart = dc.rowChart('#errors-by-type-chart');
        // var dataTable = dc.dataTable('.dc-data-table');
        // var verrorsCount = dc.dataCount('.dc-data-count');


        // var errorData;

        ref = "/validate/data/" + txID;
        d3.json(ref, function(error, data) {

            var errorsBarChart = dc.barChart("#errors-chart");
            var errorsByTypeChart = dc.rowChart('#errors-by-type-chart');
            var dataTable = dc.dataTable('.dc-data-table');
            var verrorsCount = dc.dataCount('.dc-data-count');
            var errorData = data;

            // normalize/parse data so dc can correctly sort & bin them
            errorData.forEach(function(d) {
                d.originalLine = +d.originalLine;
            });
            // console.log(errorData);

            var ndx = crossfilter(errorData);
            var all = ndx.groupAll();

            var lineDim = ndx.dimension(function(d) {
                return d.originalLine;
            });

            var typeDim = ndx.dimension(function(d) {
                return d.validationType;
            });
            var validationTypesGroup = typeDim.group();

            var allDim = ndx.dimension(function(d) {
                return d;
            });

            // var countPerLine = lineDim.group().reduceCount(function(d) {
            //     return d.originalLine;
            // });

            var lineGroup = lineDim.group();

            // var countPerLine = lineDim.group().reduceSum(function(d) {return d.OriginalLine;});

            errorsBarChart
                .width(350)
                .height(190)
                .margins({
                    top: 20,
                    right: 0,
                    bottom: 0,
                    left: 0
                })
                .gap(1)
                .x(d3.scale.linear().domain([0, 200000]))
                .elasticX(true)
                .elasticY(true)
                .dimension(lineDim)
                .group(lineGroup);

            errorsBarChart.dimension(lineDim);

            errorsByTypeChart
                .width(350)
                .height(190)
                .margins({
                    top: 20,
                    left: 10,
                    right: 10,
                    bottom: 20
                })
                .group(validationTypesGroup)
                .dimension(typeDim)
                .title(function(d) {
                    return d.value;
                })
                .elasticX(true)
                .xAxis().ticks(4);


            verrorsCount
                .dimension(ndx)
                .group(all);


            dataTable
            // .width(900)
            // .height(800)
                .dimension(allDim)
                .group(function(d) {
                    // return 'dc.js insists on putting a row here so I remove it using JS';
                    // return d.originalLine;
                    return 'Errors ordered by original file line number (table shows first 100 errors)'
                })
                .size(100)
                .columns([
                    // function (d) { return d.txID; },
                    function(d) {
                        return d.originalLine;
                    },
                    function(d) {
                        return d.validationType;
                    },
                    function(d) {
                        return d.errField;
                    },
                    function(d) {
                        return d.description;
                    }
                ])
                .sortBy(function(d) {
                    return d.originalLine;
                })
                .order(d3.ascending)
                .on('renderlet', function(table) {
                    // each time table is rendered remove nasty extra row dc.js insists on adding
                    // table.select('tr.dc-table-group').remove();
                    table.selectAll('.dc-table-group').classed('info', true);
                });

            dc.renderAll();

            // dc.redrawAll();


        });


    }
