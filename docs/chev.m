clear;
clc;
close all;

TP = [10,13,20,40,80]; %8
LatQ3 = [210.22,209.26,212.15,221.35,243]; %5
ConQ3 = [100,100,100,98.35,87.79]; %5
LatB3 = [218.79,215.65,222.64,239.23,268.78]; %5
ConB3 = [99.79,99.33,98.96,97.92,96.10];
LatDN3 = [215.61,214.66,219.84,226.79,247.30];
ConDN3 = [99.69,99.54,99.33,99.72,98.40];

T = [1000,500,250,100,75,50,25,12.5];
Con = [98.51,99.86,100,100,100,100,99.60,99.86];

Lat1R = [104.35,103.33,105.24,105.24,109.05];
Lat2R = [205.08,205.50,211.61,211.28,230.48];
Lat4R = [217.70,211.89,217.21,236.38]; %4
Lat5R = [233.14,232.18,250.10]; %3

LatCv = [321.33,333.57,350.01,365.21,426.89];

figure
hold on
plot(TP, LatQ3, 'LineWidth', 2)
plot(TP, LatB3, 'LineWidth', 2)
plot(TP, LatDN3, 'LineWidth', 2)
axis([0 90 200 280])
ax = gca;
ax.FontSize = 14; 
xlabel('Throughput (ops/s)','FontSize',14)
ylabel('Latency (ms)','FontSize',14)
legend('CmRDT: Queue','CmRDT: Naive C. Broadcast','CmRDT: Do Nothing',...
    'Location','NorthWest','FontSize',12)
grid on

figure
hold on
plot(TP, ConQ3, 'LineWidth', 2)
plot(TP, ConB3, 'LineWidth', 2)
plot(TP, ConDN3, 'LineWidth', 2)
axis([0 90 90 101])
ax = gca;
ax.FontSize = 14; 
xlabel('Throughput (ops/s)','FontSize',14)
ylabel('Consistency (%)','FontSize',14)
legend('CmRDT: Queue','CmRDT: Naive C. Broadcast','CmRDT: Do Nothing',...
    'Location','SouthWest','FontSize',12)
grid on

figure
hold on
plot(T, Con, 'k-', 'LineWidth', 2)
axis([0 1000 98 100.5])
ax = gca;
ax.FontSize = 14; 
xlabel('Queue Processing Frequency (ms)','FontSize',14)
ylabel('Consistency (%)','FontSize',14)
grid on

figure 
hold on
plot(TP, Lat1R, 'LineWidth', 2)
plot(TP, Lat2R, 'LineWidth', 2)
plot(TP, LatQ3, 'LineWidth', 2)
plot(TP(1:4), Lat4R, 'LineWidth', 2)
plot(TP(1:3), Lat5R, 'LineWidth', 2)
axis([0 90 100 260])
ax = gca;
ax.FontSize = 14; 
xlabel('Throughput (ops/s)','FontSize',14)
ylabel('Latency (ms)','FontSize',14)
legend('1 Replica','2 Replicas','3 Replicas','4 Replicas','5 Replicas',...
    'Location','EastOutside','FontSize',12)
grid on

figure 
hold on
plot(TP, LatQ3, 'LineWidth', 2)
plot(TP, LatCv, 'LineWidth', 2)
axis([0 90 200 450])
ax = gca;
ax.FontSize = 14; 
xlabel('Throughput (ops/s)','FontSize',14)
ylabel('Latency (ms)','FontSize',14)
legend('CmRDT-Queue','CvRDT',...
    'Location','NorthWest','FontSize',12)
grid on



