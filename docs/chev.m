clear;
clc;
close all;

TP = [10,25,50,75,100,150,200]; %7
LatZero = [320.3608322,321.3453256,327.0924333,326.6492689,330.8418789,333.4437322,330.8689822]; 
ConZero = [98.67052756,99.32295689,98.46180389,98.55690111,99.16398878,97.36931389,94.52633778]; 
LatCmO = [120.7411712,121.758385,131.9334411,159.2330944,204.9730367,410.9449833,1131.580602]; 
ConCmO = [100,100,100,91.25328278,74.49021689,38.28706789,30.87328556];
LatCv = [115.6390222,117.1922789,121.3292156,116.7463901,120.7267906,122.0193823,125.5836089];
ConCv = [100,99.85822889,92.95622156,92.31734,85.67499022]; %5

T = [1000,500,250,100,75,50,25,12.5];
Con = [98.51,99.86,100,100,100,100,99.60,99.86];

Lat1R = [84.219005,121.0959633,267.01762,334.73009,469.237945,640.1349267,1086.41542];
Lat5R = [150.439581,150.58367,158.1996176,179.709039,208.292013,588.74626,1639.969189];
Lat7R = [150.9895434,150.7543266,275.404859,212.2111357,210.1299633,691.4358629,2341.255343]; 

figure
hold on
plot(TP, LatZero, 'k:', 'LineWidth', 2)
plot(TP, LatCmO, 'k-.', 'LineWidth', 2)
plot(TP, LatCv, 'k-', 'LineWidth', 2)
axis([5 200 80 350])
ax = gca;
ax.FontSize = 14; 
xlabel('Throughput (ops/s)','FontSize',14)
ylabel('Latency (ms)','FontSize',14)
legend('Zero','CmRDT-O','CvRDT',...
    'Location','Best','FontSize',12)
grid on

figure
hold on
plot(TP, ConZero, 'k:', 'LineWidth', 2)
plot(TP, ConCmO, 'k-.', 'LineWidth', 2)
plot(TP(1:5), ConCv, 'k-', 'LineWidth', 2)
axis([5 200 90 101])
ax = gca;
ax.FontSize = 14; 
xlabel('Throughput (ops/s)','FontSize',14)
ylabel('Consistency (%)','FontSize',14)
legend('Zero','CmRDT-O','CvRDT',...
    'Location','Best','FontSize',12)
grid on

% figure
% hold on
% plot(T, Con, 'k-', 'LineWidth', 2)
% axis([0 1000 98 100.5])
% ax = gca;
% ax.FontSize = 14; 
% xlabel('Queue Processing Frequency (ms)','FontSize',14)
% ylabel('Consistency (%)','FontSize',14)
% grid on

figure 
hold on
plot(TP, Lat1R, 'k-', 'LineWidth', 2)
plot(TP, LatCmO, 'k:', 'LineWidth', 2)
plot(TP, Lat5R, 'k-.', 'LineWidth', 2)
plot(TP, Lat7R, 'k--', 'LineWidth', 2)
axis([5 205 80 700])
ax = gca;
ax.FontSize = 14; 
xlabel('Throughput (ops/s)','FontSize',14)
ylabel('Latency (ms)','FontSize',14)
legend('1 Replica','3 Replicas','5 Replicas','7 Replicas',...
    'Location','Best','FontSize',12)
grid on
% 
% figure 
% hold on
% plot(TP, LatQ3, 'LineWidth', 2)
% plot(TP, LatCv, 'LineWidth', 2)
% axis([0 90 200 450])
% ax = gca;
% ax.FontSize = 14; 
% xlabel('Throughput (ops/s)','FontSize',14)
% ylabel('Latency (ms)','FontSize',14)
% legend('CmRDT-Queue','CvRDT',...
%     'Location','NorthWest','FontSize',12)
% grid on
% 
% 
% 
